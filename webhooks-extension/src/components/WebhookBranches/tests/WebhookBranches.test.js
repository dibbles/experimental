/*
Copyright 2019 The Tekton Authors
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

import React from 'react';
import { waitForElement } from 'react-testing-library';
import { WebhookBranches } from '../WebhookBranches';
import * as API from '../../../../src/api/index';
import { renderWithRouter } from '../../../test/utils/test';
import 'react-testing-library/cleanup-after-each';

beforeEach(jest.restoreAllMocks);

const webhook = {
  url: "https://githuub.com/someUser/someRepo",
  namespace: "default",
  pipeline: "pipeline1"
};

const pipelineRuns = {
  items: [
    {
      metadata: {
        labels: {
          gitBranch: 'master'
        }
      },
      status: {
        conditions: [
          {
            lastTransitionTime: '2019-09-23T18:27:57Z',
            reason: 'Failed'
          }
        ]
      }
    },
    {
      metadata: {
        labels: {
          gitBranch: 'branch1'
        }
      },
      status: {
        conditions: [
          {
            lastTransitionTime: '2019-09-23T18:27:57Z',
            reason: 'Success'
          }
        ]
      }
    }
  ]
};

const fetchBranchFailMock = {
  response: {
    text: () => {
      return Promise.resolve('Mock Error Fetching Branch');
    }
  }
};

it('display branches', async () => {
  jest
    .spyOn(API, 'getPipelineRuns')
    .mockImplementation(() => Promise.resolve(pipelineRuns));

  const { getByText } = renderWithRouter(
    <WebhookBranches close={() => {}} webhook={webhook} />
  );
  await waitForElement(() => getByText(/master/i));
  await waitForElement(() => getByText(/branch1/i));
});

it('display notification when error occurs', async () => {
  jest
    .spyOn(API, 'getPipelineRuns')
    .mockImplementation(() => Promise.reject(fetchBranchFailMock));

  const { getByText } = renderWithRouter(
    <WebhookBranches close={() => {}} webhook={webhook} />
  );
  await waitForElement(() => getByText(/Mock Error Fetching Branch/i));
});
