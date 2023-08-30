// 编写一个函数，计算两个SHA256哈希码中不同bit的数目。（参考2.6.2节的PopCount函数。)
package main
import "crypto/sha256"
import "fmt"

func main() {
	b1 := []byte("x")
	b2 := []byte("X")
	c1:=sha256.Sum256(b1)
	c2:=sha256.Sum256(b2)
	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
	count := 0
	fmt.Println(c1)
	fmt.Println(c2)
	for i:=0; i<len(c1);i++ {
		if c1[i] == c2[i] {
			fmt.Println(c1[i])
			count++
		}
	}

	fmt.Print(count)
}