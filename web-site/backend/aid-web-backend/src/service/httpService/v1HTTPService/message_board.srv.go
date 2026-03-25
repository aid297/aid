package v1HTTPService

import (
	"errors"
	`fmt`
	`os`

	"github.com/gofrs/uuid/v5"
	jsonIter "github.com/json-iterator/go"

	"github.com/aid297/aid/array/anySlice"
	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
)

// MessageBoardService 服务：留言板
type MessageBoardService struct{}

func (*MessageBoardService) getDirectionFile() (
	directionFile filesystemV4.IFilesystem,
	directionFileSlice anySlice.AnySlicer[string],
	err error,
) {
	var (
		directionDir         = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.MessageBoard.Dir))
		directionFileJSON    []byte
		directionFileContent []string
	)

	if !directionDir.GetExist() {
		if err = directionDir.Create().GetError(); err != nil {
			return
		}
	}

	directionFile = filesystemV4.NewFile(filesystemV4.Rel(global.CONFIG.MessageBoard.Dir, "direction.json"))

	if !directionDir.GetExist() {
		return nil, nil, nil
	}

	if !directionFile.GetExist() {
		if err = directionFile.Write([]byte{'[', ']'}, filesystemV4.Mode(0755), filesystemV4.Flag(os.O_WRONLY|os.O_TRUNC|os.O_CREATE)).GetError(); err != nil {
			return nil, nil, fmt.Errorf("创建目录文件失败：%w", err)
		}
		return directionFile, anySlice.New[string](), nil
	}

	if directionFileJSON, err = directionFile.Read(); err != nil {
		return
	}

	if err = jsonIter.Unmarshal(directionFileJSON, &directionFileContent); err != nil {
		return
	}

	directionFileSlice = anySlice.NewList(directionFileContent)

	return
}

// List 留言板服务：获取信息列表
func (*MessageBoardService) List() (messages []map[string]string, err error) {
	var (
		directionFileSlice anySlice.AnySlicer[string]
		messageFile        filesystemV4.IFilesystem
		messageFileContent []byte
		messageContent     map[string]string
	)

	if _, directionFileSlice, err = New.MessageBoard().getDirectionFile(); err != nil {
		return
	}

	messages = make([]map[string]string, directionFileSlice.LengthNotEmpty())
	directionFileSlice.Each(func(idx int, item string) {
		messageFile = filesystemV4.NewFile(filesystemV4.Rel(global.CONFIG.MessageBoard.Dir, item))
		if !messageFile.GetExist() {
			return
		}

		messageContent = make(map[string]string)
		if messageFileContent, err = messageFile.Read(); err != nil {
			return
		}

		if err = jsonIter.Unmarshal(messageFileContent, &messageContent); err != nil {
			return
		}

		messages[idx] = messageContent
	})

	return messages, nil
}

// Store 留言板服务：保存信息
func (*MessageBoardService) Store(form *request.MessageBoardStoreRequest) (err error) {
	var (
		directionFile      filesystemV4.IFilesystem
		directionFileJSON  []byte
		directionFileSlice anySlice.AnySlicer[string]
		newUUID            = uuid.Must(uuid.NewV7())
		messageFile        filesystemV4.IFilesystem
		messageFileContent []byte
	)

	if directionFile, directionFileSlice, err = New.MessageBoard().getDirectionFile(); err != nil {
		return
	}

	directionFileSlice.Append(newUUID.String())

	if directionFileJSON, err = jsonIter.Marshal(directionFileSlice.ToSlice()); err != nil {
		return
	}

	if err = directionFile.Write(directionFileJSON, filesystemV4.Flag(os.O_CREATE|os.O_WRONLY), filesystemV4.Mode(0755)).GetError(); err != nil {
		return
	}

	messageFile = filesystemV4.NewFile(filesystemV4.Abs(directionFile.GetBasePath(), newUUID.String()))

	if messageFileContent, err = jsonIter.Marshal(map[string]string{"uuid": newUUID.String(), "content": form.Content}); err != nil {
		return
	}

	if err = messageFile.Write(messageFileContent, filesystemV4.Flag(os.O_CREATE|os.O_WRONLY), filesystemV4.Mode(0755)).GetError(); err != nil {
		return
	}

	return
}

// Destroy 留言板服务：删除信息
func (*MessageBoardService) Destroy(form *request.MessageBoardDestroyRequest) (err error) {
	var (
		directionFile      filesystemV4.IFilesystem
		directionFileJSON  []byte
		directionFileSlice anySlice.AnySlicer[string]
		idx                int
	)

	if directionFile, directionFileSlice, err = New.MessageBoard().getDirectionFile(); err != nil {
		return
	}

	if idx = directionFileSlice.GetIndexByValue(form.ID); idx == -1 {
		return errors.New("ID不存在")
	}

	directionFileSlice.RemoveByIndex(idx)

	if directionFileJSON, err = jsonIter.Marshal(directionFileSlice.ToSlice()); err != nil {
		return
	}

	if err = directionFile.Write(directionFileJSON, filesystemV4.Flag(os.O_WRONLY|os.O_TRUNC|os.O_CREATE), filesystemV4.Mode(0755)).GetError(); err != nil {
		return
	}

	messageFile := filesystemV4.NewFile(filesystemV4.Abs(directionFile.GetBasePath(), form.ID))
	if err = messageFile.Remove().GetError(); err != nil {
		return
	}

	return
}
