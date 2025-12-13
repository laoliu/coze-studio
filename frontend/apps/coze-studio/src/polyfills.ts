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

/**
 * Polyfills for browser environment
 * 注入 process 和其他 Node.js 全局变量
 */

import process from 'process/browser.js';

// 注入 process 到全局
if (typeof window !== 'undefined') {
  (window as any).process = process;
  (window as any).global = window;
}

// 设置默认的 NODE_ENV
if (!process.env.NODE_ENV) {
  process.env.NODE_ENV = 'production';
}

export {};
