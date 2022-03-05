package rope

type Settings struct {
	SplitLength int     // Maximum length before to split a rope
	JoinLength  int     // Minimum length to join a rope
	Rebalance   float32 // Ratio needed to rebalance a rope
}

var DefaultSettings = &Settings {
	SplitLength: 400,
	JoinLength:  200,
	Rebalance:   1.5,
}

type Rope[T any] struct {
	value    []T
	length   int
	left     *Rope[T]
	right    *Rope[T]
	settings *Settings
}

func NewRope[T any](value []T, settings *Settings) *Rope[T] {
	rope := &Rope[T]{value: value, length: len(value), settings: settings}
	rope.adjust()
	return rope
}

func (r *Rope[T]) adjust() {
	if r.value != nil && r.length > r.settings.SplitLength { // It is not yet split but too long
		r.left  = NewRope[T](r.value[:r.length / 2], r.settings)
		r.right = NewRope[T](r.value[r.length / 2:], r.settings)
		r.value = nil // Mark as split
		return
	}
	if r.value == nil && r.length < r.settings.JoinLength { // It is split but too short
		r.value = make([]T, r.length)
		r.left.Copy(r.value)
		r.right.Copy(r.value[:r.left.length])
		r.left = nil
		r.right = nil
	}
}

func (r *Rope[T]) Remove(start, end int) *Rope[T] {
	if start == end {
		return r
	}
	if r.value != nil { // If rope isn't split
		// A copy is needed, as append doesn't guarantee immutability
		newValue := make([]T, r.length - (end - start))
		copy(newValue, r.value[:start])
		copy(newValue[start:], r.value[end:])
		changed := NewRope(newValue, r.settings)
		return changed
	}
	// Rope is split
	changed := &Rope[T]{settings: r.settings}
	leftStart, leftEnd := bound(start, end, r.left.length)
	changed.left = r.left.Remove(leftStart, leftEnd)

	rightStart, rightEnd := bound(start - r.left.length, end - r.left.length, r.right.length)
	changed.right = r.right.Remove(rightStart, rightEnd)

	changed.length = changed.left.length + changed.right.length
	changed.adjust()
	return changed
}

func (r *Rope[T]) Insert(index int, insertion []T) *Rope[T] {
	if r.value != nil { // If rope isn't split
		// A copy is needed, as append doesn't guarantee immutability
		newValue := make([]T, r.length + len(insertion))
		copy(newValue, r.value[:index])
		copy(newValue[index:], insertion)
		copy(newValue[index + len(insertion):], r.value[index:])
		changed := NewRope(newValue, r.settings) // Takes care of adjusting
		return changed
	}
	// Rope is split
	changed := &Rope[T]{
		settings: r.settings,
		length: r.length + len(insertion),
		left: r.left,
		right: r.right,
	}

	if index < r.left.length {
		changed.left = r.left.Insert(index, insertion)
	} else {
		changed.right = r.right.Insert(index - r.left.length, insertion)
	}
	return changed
}

// Bind the start and end indexes inside a length,
// preventing OOB.
func bound(start, end, length int) (newStart, newEnd int) {
	if start < 0 {
		start = 0
	} else if start > length {
		start = length
	}
	if end < 0 {
		end = 0
	} else if end > length {
		end = length
	}
	return start, end
}

func (r *Rope[T]) Copy(dst []T) {
	if r.value != nil {
		copy(dst, r.value)
	} else {
		r.left.Copy(dst)
		r.right.Copy(dst[r.left.length:])
	}
}

func (r *Rope[T]) CopySlice(dst []T, start, end int) {
	if start == end {
		return
	}
	if r.value != nil { // Isn't split
		copy(dst, r.value[start:end])
		return
	}
	// Is split
	leftStart, leftEnd := bound(start, end, r.left.length)
	r.left.CopySlice(dst, leftStart, leftEnd)

	rightStart, rightEnd := bound(start - r.left.length, end - r.left.length, r.right.length)
	r.right.CopySlice(dst[leftEnd - leftStart:], rightStart, rightEnd)
}

func (r *Rope[T]) Value() []T {
	value := make([]T, r.length)
	r.Copy(value)
	return value
}

func (r *Rope[T]) Slice(start, end int) []T {
	value := make([]T, end - start)
	r.CopySlice(value, start, end)
	return value
}

func (r *Rope[T]) Length() int {
	return r.length
}

// NOTE: This is a very slow way to do things
func (r *Rope[T]) Rebalance() {
	if r.value != nil {
		return
	}
	if float32(r.left.length) / float32(r.right.length) > r.settings.Rebalance ||
	   float32(r.right.length) / float32(r.left.length) > r.settings.Rebalance {
		   rebalancedRope := NewRope(r.Value(), r.settings)
		   *r = *rebalancedRope
	} else {
		r.left.Rebalance()
		r.right.Rebalance()
	}
}
