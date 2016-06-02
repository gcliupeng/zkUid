package zkUid
import (
	"time"
	"strings"
	"strconv"
	//"path"
	"github.com/wandoulabs/go-zookeeper/zk"
	"fmt"
)
type ZkUidInstance struct{
	Host string
	Path string
	Conn *zk.Conn
}
func NewInstance(zkhost,zkpath string) (*ZkUidInstance,error) {
	c, _, err := zk.Connect([]string{zkhost}, time.Second, 30) //*10)
	if err != nil {
		return nil,err
	}
	return &ZkUidInstance{zkhost,zkpath,c},nil
}
func (zkUid *ZkUidInstance) GetUid() (int,error) {
loop:	
	children, _, err := zkUid.Conn.Children(zkUid.Path)
		if err != nil {
			//fmt.Printf("%v",children)
			fmt.Printf("%v",err)
			return -1,err
		}
	subPath:=""
	if len(children)==0 {
		subPath,err=zkUid.Conn.Create(zkUid.Path+"/", []byte(""),zk.FlagSequence, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Printf("%v",err)
			return -1,err
		}
		partsT:=strings.Split(subPath, "/")
		subPath=partsT[len(partsT)-1]
	}else{
		max:=-1
		for _,p :=range children{
			tInt,_:=strconv.Atoi(p)
			if tInt>max{
				max=tInt
				subPath=p
			}
		}
	}

	tmpPath,err:=zkUid.Conn.Create(zkUid.Path+"/"+subPath+"/", []byte(""),zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Printf("%v",err)
		return -1,err
		}
	parts:=strings.Split(tmpPath, "/")
	tmpPath=parts[len(parts)-1]
	big,err:=strconv.Atoi(tmpPath)
	//zk sequence 节点最大为2147483647
	if big==2147483647 {
		zkUid.Conn.Create(zkUid.Path+"/", []byte(""),zk.FlagSequence, zk.WorldACL(zk.PermAll))
		goto loop
	}else{
		subPathI,_:=strconv.Atoi(subPath)
		tmpPathI,_:=strconv.Atoi(tmpPath)
		//fmt.Printf("%d\n",subPathI*2147483647+tmpPathI)
		//fmt.Printf("%v\n",subPath)
		//fmt.Printf("%v\n",tmpPath)
		return subPathI*2147483647+tmpPathI,nil
	}

}