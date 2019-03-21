package utils

import "sort"


//判断两个int slice包含的值是否相等，顺序不一致内容一致也是相等
func IsEqualIntSlice(a []int, b []int) bool {
    sort.Ints(a)
    sort.Ints(b)

    if len(a) != len(b) {
        return false
    }

    for i, v := range a {
        if v != b[i] {
            return false
        }
    }

    return true
}

func InIntSlice(e int, a []int) bool {
    for _, v := range a {
        if v == e {
            return true
        }
    }
    return false
}

func InSlice(e string, a []string) bool {
    for _, v := range a {
        if v == e {
            return true
        }
    }
    return false
}

func SliceDiff(a, b []string) (diff []string) {
    for _, v := range a {
        if !InSlice(v, b) {
            diff = append(diff, v)
        }
    }
    return
}

func SliceIntersect(a, b []string) (intersec []string) {
    for _, v := range a {
        if InSlice(v, b) {
            intersec = append(intersec, v)
        }
    }
    return
}
