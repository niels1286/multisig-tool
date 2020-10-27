/*
Copyright © 2020 Niels niels@nuls.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/niels1286/multisig-tool/cmd"
	"github.com/niels1286/multisig-tool/i18n"
	"time"
)

func main() {
	i18n.InitLang(getLangType())
	cmd.Execute()
}

func getLangType() string {
	zone, offset := time.Now().Local().Zone()
	if zone == "CST" && offset == 28800 {
		return "cn"
	}
	return "en"
}
