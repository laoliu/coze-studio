/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

// PluginInfo 插件信息
type PluginInfo struct {
	Manifest  *Manifest `json:"manifest"`
	ServerURL *string   `json:"server_url,omitempty"`
}

// Manifest 插件清单
type Manifest struct {
	NameForModel string `json:"name_for_model"`
	NameForHuman string `json:"name_for_human"`
	Description  string `json:"description"`
}
