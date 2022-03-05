# Rope
A general Go implementation of a [persistent rope data structure](<https://en.wikipedia.org/wiki/Rope_%28data_structure%29>),
using the new 1.18 generics.
It uses only nondestructive operations,
so taking snapshots has no extra cost.

Using the default settings, its faster than a monolithic array since
between 100 and 1000 elements for random insertions.

The implementation was inspired by [this javascript repo](https://github.com/component/rope),
which is surprinsingly readable.

## Usage
```go
func main() {
	firstRope := rope.NewRope([]float32{0, 1, 2, 3}, rope.DefaultSettings) // Can be any type
	
	secondRope := firstRope.Remove(1, 2)
	thirdRope := secondRope.Insert(1, []float32{1.2, 1.5, 1.9})

	fmt.Println(firstRope.Value())  // [0, 1, 2, 3]
	fmt.Println(secondRope.Value()) // [0, 2, 3]
	fmt.Println(thirdRope.Value())  // [0, 1.2, 1.5, 1.9, 2, 3]


	thirdRope.Rebalace() // Ensures the rope is properly balanced, doesn't change value
	fmt.Println(thirdRope.Value())  // [0, 1.2, 1.5, 1.9, 2, 3]

	fmt.Println(thirdRope.Slice(1, 3))  // [1.2, 1.5]

	value := make([]float32, thirdRope.Length())
	thirdRope.Copy(value)
	fmt.Println(value) // [0, 1.2, 1.5, 1.9, 2, 3]

	slice := make([]float32, 2)
	thirdRope.CopySlice(value, 1, 3)
	fmt.Println(slice) // [1.2, 1.5]
}
```
