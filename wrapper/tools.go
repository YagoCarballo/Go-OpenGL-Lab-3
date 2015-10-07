package wrapper
import "io/ioutil"

func ReadFile (path string) (string, error) {
	fileContents, err := ioutil.ReadFile(path)
	return string(fileContents) + "\x00", err
}

func FileToString (path string) string {
	content, err := ReadFile(path)
	if err != nil {
		panic(err)
	}

	return content
}