package liquid

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/karlseguin/liquid"
)

func Parse(path string) ([]byte, error) {
	content, err := os.ReadFile("C:/Users/uinme/go_workspace/etl/t_goetl_test.yml.liquid")
	if err != nil {
		return nil, err
	}

	config := liquid.Configure().IncludeHandler(includeHandler)

	template, err := liquid.ParseString(string(content), config)
	writer := &bytes.Buffer{}
	// todo 環境変数の値を取得して第2引数に渡す
	template.Render(writer, nil)

	return writer.Bytes(), nil
}

func includeHandler(name string, writer io.Writer, data map[string]interface{}) {
	config := liquid.Configure().IncludeHandler(includeHandler)
	// todo 入力ファイルのディレクトリと同じ階層のものを読み込むように修正
	fileName := path.Join("./", "_"+name+".yml.liquid")
	template, _ := liquid.ParseFile(fileName, config)
	template.Render(writer, data)
}
