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
	configDir := path.Join(userHomeDir, ".vptun")
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

	viper.SetEnvPrefix("vptun-certificate")
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

	outputServerHttpCACrt := path.Join(output, "server/http/ca.crt")
	outputServerHttpCAKey := path.Join(output, "server/http/ca.key")
	outputServerHttpServerCrt := path.Join(output, "server/http/server.crt")
	outputServerHttpServerKey := path.Join(output, "server/http/server.key")

	outputServerGrpcCACrt := path.Join(output, "server/grpc/ca.crt")
	outputServerGrpcCAKey := path.Join(output, "server/grpc/ca.key")
	outputServerGrpcServerCrt := path.Join(output, "server/grpc/server.crt")
	outputServerGrpcServerKey := path.Join(output, "server/grpc/server.key")
	outputServerGrpcClientCrt := path.Join(output, "server/grpc/client.crt")
	outputServerGrpcClientKey := path.Join(output, "server/grpc/client.key")

	outputClientGrpcCaCrt := path.Join(output, "client/grpc/ca.crt")
	outputClientGrpcClientCrt := path.Join(output, "client/grpc/client.crt")
	outputClientGrpcClientKey := path.Join(output, "client/grpc/client.key")

	var errs []error
	if err = fileExists(outputServerHttpCACrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerHttpCAKey); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerHttpServerCrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerHttpServerKey); err != nil {
		errs = append(errs, err)
	}

	if err = fileExists(outputServerGrpcCACrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerGrpcCAKey); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerGrpcServerCrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerGrpcServerKey); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerGrpcClientCrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputServerGrpcClientKey); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputClientGrpcCaCrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputClientGrpcClientCrt); err != nil {
		errs = append(errs, err)
	}
	if err = fileExists(outputClientGrpcClientKey); err != nil {
		errs = append(errs, err)
	}
	if len(errs) != 0 {
		for _, err = range errs {
			fmt.Println(err)
		}
		err = nil
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

	os.MkdirAll(path.Join(output, "server/http"), 0755)
	os.MkdirAll(path.Join(output, "server/grpc"), 0755)
	os.MkdirAll(path.Join(output, "client/grpc"), 0755)

	// 写入
	if err = ioutil.WriteFile(outputServerHttpCAKey, []byte(httpCA.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerHttpCACrt, []byte(httpCA.Certificate), 0600); err != nil {
		return
	}

	if err = ioutil.WriteFile(outputServerHttpServerKey, []byte(httpServer.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerHttpServerCrt, []byte(httpServer.Certificate), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerGrpcCAKey, []byte(grpcCA.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerGrpcCACrt, []byte(grpcCA.Certificate), 0600); err != nil {
		return
	}

	if err = ioutil.WriteFile(outputServerGrpcServerKey, []byte(grpcServer.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerGrpcServerCrt, []byte(grpcServer.Certificate), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerGrpcClientKey, []byte(grpcClient.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputServerGrpcClientCrt, []byte(grpcClient.Certificate), 0600); err != nil {
		return
	}

	// client
	if err = ioutil.WriteFile(outputClientGrpcCaCrt, []byte(grpcCA.Certificate), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputClientGrpcClientKey, []byte(grpcClient.PrivateKey), 0600); err != nil {
		return
	}
	if err = ioutil.WriteFile(outputClientGrpcClientCrt, []byte(grpcClient.Certificate), 0600); err != nil {
		return
	}

	fmt.Println("")
	fmt.Println("Server Http CA PrivateKey : ", outputServerHttpCAKey)
	fmt.Println("Server Http CA Certificate: ", outputServerHttpCACrt)
	fmt.Println("")
	fmt.Println("Server Http Server PrivateKey : ", outputServerHttpServerKey)
	fmt.Println("Server Http Server Certificate: ", outputServerHttpServerCrt)
	fmt.Println("")
	fmt.Println("Server Grpc CA PrivateKey : ", outputServerGrpcCAKey)
	fmt.Println("Server Grpc CA Certificate: ", outputServerGrpcCACrt)
	fmt.Println("")
	fmt.Println("Server Grpc Server PrivateKey : ", outputServerGrpcServerKey)
	fmt.Println("Server Grpc Server Certificate: ", outputServerGrpcServerCrt)
	fmt.Println("")
	fmt.Println("Server Grpc Client PrivateKey : ", outputServerGrpcClientKey)
	fmt.Println("Server Grpc Client Certificate: ", outputServerGrpcClientCrt)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Client Grpc CA Certificate    : ", outputClientGrpcCaCrt)
	fmt.Println("Client Grpc Client PrivateKey : ", outputClientGrpcClientKey)
	fmt.Println("Client Grpc Client Certificate: ", outputClientGrpcClientCrt)
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
