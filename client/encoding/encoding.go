package encoding

import(
	"fmt"
	"github.com/manpreet-singh1040/ByteBeam/client/fileHandling"
)
func Encoding(){
	fmt.Println("hello from encoding package")
	address := "D:\\series\\friends\\season 1\\Firends.S01.720p.English.Vegamovies.NL\\Friends.S01E19.720p.English.Vegamovies.NL.mkv"
	file := fileHandling.File{}
	file.Initialize(address)
	// fmt.Println("file metadata is --> ",file.Metadata)
}
