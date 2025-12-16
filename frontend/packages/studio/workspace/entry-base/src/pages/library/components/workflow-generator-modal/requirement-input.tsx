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

import React from 'react';
import { TextArea } from '@coze-arch/coze-design';

interface RequirementInputProps {
  requirement: string;
  setRequirement: (value: string) => void;
}

export const RequirementInput: React.FC<RequirementInputProps> = ({
  requirement,
  setRequirement,
}) => {
  return (
    <div style={{ marginBottom: '16px' }}>
      <div style={{ marginBottom: '8px', fontWeight: 600, fontSize: '14px' }}>
        Workflow Requirements <span style={{ color: '#ff4d4f' }}>*</span>
      </div>
      <TextArea
        value={requirement}
        onChange={setRequirement}
        placeholder="Describe your workflow requirements, e.g.: I need a workflow to analyze user feedback and automatically categorize it"
        rows={4}
        maxLength={1000}
      />
    </div>
  );
};
