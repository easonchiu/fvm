package script

import (
	"fmt"
	"fvm/constant"
	"fvm/utils"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// 安装指定版本
func InstallVersion(version string) error {
	json, err := utils.GetYarnPackageJson(constant.FEC_BUILDER)
	if err != nil {
		return err
	}

	// 在 versions 列表中查询指定版本
	versions := json.Get("versions").Map()
	if _, ok := versions[version]; !ok {
		return fmt.Errorf("the \"%v\" version could not be found", version)
	}

	// 查询 tarball，这是对应版本的 .tgz 格式的下载地址
	tarball := versions[version].Get("dist").Get("tarball").String()
	if len(tarball) == 0 {
		return fmt.Errorf("the \"%v\" version could not be found", version)
	}

	// 保存到该位置
	fvmDir := utils.GetFvmDir()
	localPath := path.Join(fvmDir, constant.FEC_BUILDER+"_"+version+".tgz")
	localPathWithoutTgz := strings.ReplaceAll(localPath, ".tgz", "")

	// 判断 .tgz 是否存在
	exists := utils.CheckFileExists(localPath)
	if exists {
		return fmt.Errorf("you have already installed version \"%v\"", version)
	}

	err = utils.DownloadPkgTgz(tarball, localPath)
	if err != nil {
		return err
	}

	// 解压文件
	err = utils.DecompressTgz(localPath, localPathWithoutTgz)
	if err != nil {
		return err
	}

	return nil
}

// 切换到指定版本
func SwitchVersion(version string) error {
	list := GetLocalVersionList()

	// 判断本地是有否该版本
	exists := false
	for _, v := range list {
		if v == version {
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("the version \"%v\" was not found\nplease use [fvm install %v] to install it", version, version)
	}

	return nil
}

// 返回本地的已安装版本列表
func GetLocalVersionList() []string {
	fvmDir := utils.GetFvmDir()
	list := make([]string, 0)

	entries, err := os.ReadDir(fvmDir)
	if err != nil {
		return list
	}

	// 匹配目录名称，理论上目录名称都是 fec-builder_2.7.1 这样子的格式的
	pathReg := regexp.MustCompile("^" + constant.FEC_BUILDER + "_[0-9.]+$")
	filePrefix := constant.FEC_BUILDER + "_"

	for _, entry := range entries {
		if entry.IsDir() && pathReg.MatchString(entry.Name()) {
			list = append(list, strings.ReplaceAll(entry.Name(), filePrefix, ""))
		}
	}

	return list
}

// 返回当前正在使用的版本
func GetCurrentVersion() (string, error) {
	p, err := utils.GetNpmRootPath()
	if err != nil {
		return "", err
	}

	// 查询 npm 全局中的 fec-builder 版本
	bytes, err := os.ReadFile(path.Join(p, constant.FEC_BUILDER, "package.json"))
	if err != nil {
		return "", err
	}

	// 获取 version 字段的值
	result := gjson.GetBytes(bytes, "version")

	return result.String(), nil
}
