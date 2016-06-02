package main
import(
	"fmt"
	zkuid "zkUid/lib"
)
func main() {
	zkUidInstance,err:=zkuid.NewInstance("127.0.0.1:2182","/zk/zkUid")
	if err!=nil {
		fmt.Printf("%v",err)
		return
	}
	uid,err:=zkUidInstance.GetUid()
	if err!=nil {
		fmt.Printf("%v",err)
		return
	}
	fmt.Printf("%d\n",uid)
}