package utils;

type IterableSlice struct {
    x int;
    s [][]float64;
};

func (s *IterableSlice) Next() (value []float64) {
    s.x++;
    return s.s[s.x];
};

func (s *IterableSlice) HasNext() (ok bool) {
    s.x++;
    if s.x >= len(s.s) {
        s.x--;
        return false;
    }
    s.x--;
    return true;
};

func (s *IterableSlice) GetS() [][]float64 {
    return s.s;
};

func NewIterator(s [][]float64) *IterableSlice {
    return &IterableSlice{-1, s};
};
