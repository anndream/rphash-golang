package stream;

import (
    "math/rand"
    "github.com/wenkesj/rphash/utils"
    "github.com/wenkesj/rphash/types"
    "github.com/wenkesj/rphash/defaults"
);

type Stream struct {
    runCount int;
    counts []int64;
    centroids [][]float64;
    variance float64;
    centroidCounter types.CentroidItemSet;
    randomSeedGenerator *rand.Rand;
    rphashObject types.RPHashObject;
    lshGroup []types.LSH;
    decoder types.Decoder;
    projector types.Projector;
    hash types.Hash;
    varianceTracker types.StatTest;
};

func NewStream(rphashObject types.RPHashObject) *Stream {
    randomSeedGenerator := rand.New(rand.NewSource(rphashObject.GetRandomSeed()));
    hash := defaults.NewHash(rphashObject.GetHashModulus());
    decoder := rphashObject.GetDecoderType();
    varianceTracker := defaults.NewStatTest(0.01);
    projections := rphashObject.GetNumberOfProjections();
    k := rphashObject.GetK() * projections;
    centroidCounter := defaults.NewCentroidCounter(k);
    lshGroup := make([]types.LSH, projections);
    var projector types.Projector;
    for i := 0; i < projections; i++ {
        projector = defaults.NewProjector(rphashObject.GetDimensions(), decoder.GetDimensionality(), randomSeedGenerator.Int63());
        lshGroup[i] = defaults.NewLSH(hash, decoder, projector);
    }
    return &Stream{
        counts: nil,
        centroids: nil,
        variance: 0,
        runCount: 0,
        centroidCounter: centroidCounter,
        randomSeedGenerator: randomSeedGenerator,
        rphashObject: rphashObject,
        lshGroup: lshGroup,
        hash: hash,
        decoder: decoder,
        projector: projector,
        varianceTracker: varianceTracker,
    };
};

func (this *Stream) AddVectorOnlineStep(vec []float64) int64 {
    var hash []int64;
    c := defaults.NewCentroidStream(vec);

    tmpvar := this.varianceTracker.UpdateVarianceSample(vec);

    if this.variance != tmpvar {
        for _, lsh := range this.lshGroup {
            lsh.UpdateDecoderVariance(tmpvar);
        }
        this.variance = tmpvar;
    }

    for _, lsh := range this.lshGroup {
        hash = lsh.LSHHashStream(vec, this.rphashObject.GetNumberOfBlurs());

        for _, h := range hash {
            c.AddID(h);
            // c.Centroid();
        }
    }

    this.centroidCounter.Add(c);
    return this.centroidCounter.GetCount();
};

func (this *Stream) GetCentroids() [][]float64 {
    if this.centroids == nil {
        if this.runCount == 0 {
          this.Run();
        }
        var centroids [][]float64;
        for _, cent := range this.centroidCounter.GetTop() {
            centroids = append(centroids, cent.Centroid());
        }
        this.centroids = defaults.NewKMeansStream(this.rphashObject.GetK(), centroids, this.centroidCounter.GetCounts()).GetCentroids();
    }
    return this.centroids;
};

func (this *Stream) GetVectors() [][]float64 {
  return this.rphashObject.GetVectorIterator().GetS();
};

func (this *Stream) AppendVector(vector []float64) {
    this.rphashObject.AppendVector(vector);
};

func (this *Stream) GetCentroidsOfflineStep() [][]float64 {
    var centroids [][]float64;
    var counts []int64;
    for i := 0; i < len(this.centroidCounter.GetTop()); i++ {
        centroids = append(centroids, this.centroidCounter.GetTop()[i].Centroid());
        counts = append(counts, this.centroidCounter.GetCounts()[i]);
    }
    this.centroids = defaults.NewKMeansStream(this.rphashObject.GetK(), centroids, counts).GetCentroids();
    count := int((utils.Max(counts) + utils.Min(counts)) / 2);
    counts = []int64{};
    for i := 0; i < this.rphashObject.GetK(); i++ {
        counts = append(counts, int64(count));
    }
    this.counts = counts;
    return this.centroids;
};

func (this *Stream) Run() {
    this.runCount++;
    vecs := this.rphashObject.GetVectorIterator();
    for vecs.HasNext() {
        this.AddVectorOnlineStep(vecs.Next());
    }
};
