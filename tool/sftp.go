package tool

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	user = "root"
	// 程序需在root用户下执行（or opsae）
	privateKeyPath = "/root/.ssh/id_rsa"
)

// 目标机器
type RemoteHost struct {
	Host       string
	User       string
	PrivateKey []byte
	sftpClient *sftp.Client
}

// new
func NewRemoteHost(host string) (*RemoteHost, error) {
	// 这里读取的key包含PEM编码私钥，需要ssh.ParsePrivateKey对读取到的字节数据进行解析
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &RemoteHost{
		Host:       host,
		User:       user,
		PrivateKey: key,
		// 面向对象单一职责原则(SRP)，所以*sftp.Client的初始化就放在func (h *RemoteHost) Conn() error{......}实现
	}, nil
}

// sftp.Client初始化
func (h *RemoteHost) Conn() error {
	// 解析读取的key
	singer, err := ssh.ParsePrivateKey(h.PrivateKey)
	if err != nil {
		return err
	}

	// ssh conf
	sshConf := &ssh.ClientConfig{
		User:            h.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(singer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// ssh conn
	sshConn, err := ssh.Dial("tcp", h.Host+":22", sshConf)
	if err != nil {
		// fmt.Printf("ssh conn failed! err:%v\n", err)
		return err
	}

	// sftp client
	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		return err
	}

	h.sftpClient = sftpClient

	return nil
}

// 传输文件: 相同目录下，支持一次传输多个文件
// 所有循环结束后才（defer）被动释放文件句柄，适用于较少的文件拷贝！！！
func (h *RemoteHost) CopyFile(localFilePath, remoteFilePath string, localFileName []string) error {
	for _, v := range localFileName {
		// 组装abs path
		localFile := filepath.Join(localFilePath, v)
		remoteFile := filepath.Join(remoteFilePath, v)

		// 打开local_f
		local_f, err := os.Open(localFile)
		if err != nil {
			return err
		}
		defer local_f.Close()

		// 创建/覆盖 远程f
		remote_f, err := h.sftpClient.Create(remoteFile)
		if err != nil {
			return err
		}
		defer remote_f.Close()

		// copy
		_, err = io.Copy(remote_f, local_f)
		if err != nil {
			return err
		}
	}

	return nil
}

// 传输文件: 相同目录下，支持一次传输多个文件
// 每次循环结束后就主动释放文件句柄，适用于超级多的文件拷贝！！！
func (h *RemoteHost) CopyFileMulti(localFilePath, remoteFilePath string, localFileName []string) error {
	for _, v := range localFileName {
		// fmt.Printf("start copy file: %v\n", v)
		// 组装abs path
		localFile := filepath.Join(localFilePath, v)
		remoteFile := filepath.Join(remoteFilePath, v)

		// 打开local_f
		local_f, err := os.Open(localFile)
		if err != nil {
			return err
		}

		// 创建/覆盖 远程f
		remote_f, err := h.sftpClient.Create(remoteFile)
		if err != nil {
			// 打开remote_f失败，被动关闭已经打开的local_f！
			local_f.Close()
			return err
		}

		// copy
		_, err = io.Copy(remote_f, local_f)
		if err != nil {
			// 拷贝发生错误，被动异常关闭已经打开的local_f和remote_f
			local_f.Close()
			remote_f.Close()
			return err
		}

		// 拷贝成功，主动关闭local_f
		if err = local_f.Close(); err != nil {
			// 主动关闭local_f错误，退出前主动关闭未关闭的remote_f
			remote_f.Close()
			return err
		}

		// 拷贝成功，主动关闭local_f成功，再主动关闭remote_f
		if err = remote_f.Close(); err != nil {
			return nil
		}
		// fmt.Printf("copy file : %v ok!\n", v)

	}

	return nil
}

// 释放sftp客户端
func (h *RemoteHost) Close() error {
	if h.sftpClient == nil {
		return nil
	}

	return h.sftpClient.Close()
}
