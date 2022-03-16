package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	libcertificate "github.com/otamoe/go-library/http/certificate"
	libviper "github.com/otamoe/go-library/viper"
	"github.com/spf13/viper"
)

func init() {

	// 类型
	libviper.SetDefault("http.type", "rsa", "Type of certificate issued by http")

	// 长度
	libviper.SetDefault("http.bits", 4096, " Bits of certificate issued by http")

	// host
	libviper.SetDefault("http.host", "localhost", "Host of certificate issued by http")

	// 类型
	libviper.SetDefault("grpc.type", "rsa", "Type of certificate issued by grpc")

	// 长度
	libviper.SetDefault("grpc.bits", 4096, " Bits of certificate issued by grpc")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := path.Join(userHomeDir, ".ntun")
	// 输出文件名
	libviper.SetDefault("output", configDir, "output dir")

}

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

	viper.SetEnvPrefix("ntun-certificate")
	if err = libviper.Parse(); err != nil {
		return
	}

	if libviper.PrintDefaults() {
		return
	}

	output := viper.GetString("output")

	if output == "" {
		if output, err = os.Getwd(); err != nil {
			return
		}
	}

	outputHttpCACrt := path.Join(output, "http/ca.crt")
	outputHttpCAKey := path.Join(output, "http/ca.key")
	outputHttpServerCrt := path.Join(output, "http/server.crt")
	outputHttpServerKey := path.Join(output, "http/server.key")

	outputGrpcCACrt := path.Join(output, "grpc/ca.crt")
	outputGrpcCAKey := path.Join(output, "grpc/ca.key")
	outputGrpcServerCrt := path.Join(output, "grpc/server.crt")
	outputGrpcServerKey := path.Join(output, "grpc/server.key")
	outputGrpcClientCrt := path.Join(output, "grpc/client.crt")
	outputGrpcClientKey := path.Join(output, "grpc/client.key")

	if err = fileExists(outputHttpCACrt); err != nil {
		return
	}
	if err = fileExists(outputHttpCAKey); err != nil {
		return
	}
	if err = fileExists(outputHttpServerCrt); err != nil {
		return
	}
	if err = fileExists(outputHttpServerKey); err != nil {
		return
	}

	if err = fileExists(outputGrpcCACrt); err != nil {
		return
	}
	if err = fileExists(outputGrpcCAKey); err != nil {
		return
	}
	if err = fileExists(outputGrpcServerCrt); err != nil {
		return
	}
	if err = fileExists(outputGrpcServerKey); err != nil {
		return
	}
	if err = fileExists(outputGrpcClientCrt); err != nil {
		return
	}
	if err = fileExists(outputGrpcClientKey); err != nil {
		return
	}

	httpType := viper.GetString("http.type")
	httpBits := viper.GetInt("http.bits")
	httpHost := viper.GetString("http.host")

	grpcType := viper.GetString("grpc.type")
	grpcBits := viper.GetInt("grpc.bits")

	var httpCA *libcertificate.Certificate
	var httpServer *libcertificate.Certificate

	var grpcCA *libcertificate.Certificate
	var grpcServer *libcertificate.Certificate
	var grpcClient *libcertificate.Certificate

	if httpCA, err = libcertificate.CreateTLSCertificate(httpType, httpBits, "Root http ca", nil, true, nil); err != nil {
		return
	}
	if httpServer, err = libcertificate.CreateTLSCertificate(httpType, httpBits, httpHost, []string{httpHost}, false, httpCA); err != nil {
		return
	}
	if grpcCA, err = libcertificate.CreateTLSCertificate(grpcType, grpcBits, "Root grpc ca", nil, true, nil); err != nil {
		return
	}
	if grpcServer, err = libcertificate.CreateTLSCertificate(grpcType, grpcBits, "server", []string{"server"}, false, grpcCA); err != nil {
		return
	}
	if grpcClient, err = libcertificate.CreateTLSCertificate(grpcType, grpcBits, "client", []string{"client"}, false, grpcCA); err != nil {
		return
	}

	os.MkdirAll(path.Join(output, "http"), 0755)
	os.MkdirAll(path.Join(output, "grpc"), 0755)

	// 写入
	if err = ioutil.WriteFile(outputHttpCAKey, []byte(httpCA.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputHttpCACrt, []byte(httpCA.Certificate), 0600); err != nil {
		return
	}

	if err = ioutil.WriteFile(outputHttpServerKey, []byte(httpServer.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputHttpServerCrt, []byte(httpServer.Certificate), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputGrpcCAKey, []byte(grpcCA.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputGrpcCACrt, []byte(grpcCA.Certificate), 0600); err != nil {
		return
	}

	if err = ioutil.WriteFile(outputGrpcServerKey, []byte(grpcServer.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputGrpcServerCrt, []byte(grpcServer.Certificate), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputGrpcClientKey, []byte(grpcClient.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputGrpcClientCrt, []byte(grpcClient.Certificate), 0600); err != nil {
		return
	}

	fmt.Println("")
	fmt.Println("Http CA PrivateKey : ", outputHttpCAKey)
	fmt.Println("Http CA Certificate: ", outputHttpCACrt)
	fmt.Println("")
	fmt.Println("Http Server PrivateKey : ", outputHttpServerKey)
	fmt.Println("Http Server Certificate: ", outputHttpServerCrt)
	fmt.Println("")
	fmt.Println("Grpc CA PrivateKey : ", outputGrpcCAKey)
	fmt.Println("Grpc CA Certificate: ", outputGrpcCACrt)
	fmt.Println("")
	fmt.Println("Grpc Server PrivateKey : ", outputGrpcServerKey)
	fmt.Println("Grpc Server Certificate: ", outputGrpcServerCrt)
	fmt.Println("")
	fmt.Println("Grpc Client PrivateKey : ", outputGrpcClientKey)
	fmt.Println("Grpc Client Certificate: ", outputGrpcClientCrt)
	fmt.Println("")
}

func fileExists(p string) (err error) {
	// 文件已存在
	if _, err = os.Stat(p); err == nil {
		err = fmt.Errorf("File %s already exists", p)
		return
	}

	// 其他错误
	if !os.IsNotExist(err) {
		return
	}
	err = nil
	return
}
