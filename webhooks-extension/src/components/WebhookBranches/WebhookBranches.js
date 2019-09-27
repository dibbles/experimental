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

import React, { Component } from 'react';
import { Modal } from 'carbon-components-react';
import { getPipelineRuns } from '../../api';

import './WebhookBranches.scss';

import {
  DataTable,
  DataTableSkeleton,
  InlineNotification
} from 'carbon-components-react';

const {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableBody,
  TableCell,
  TableHeader
} = DataTable;

export class WebhookBranches extends Component {
  constructor(props) {
    super(props);
    this.state = {
      rows: [],
      loading: true,
      error: null
    };
  }

  componentDidMount() {
    let { url, namespace, pipeline } = this.props.webhook;
    let [server, org, repo] = url
      .toLowerCase()
      .replace(/https?:\/\//, "")
      .split("/");

    getPipelineRuns({
      namespace,
      pipeline,
      filters: [`gitOrg=${org}`, `gitServer=${server}`, `gitRepo=${repo}`]
    })
      .then(pipelineRuns => {
        let branches = [];
        const rows = pipelineRuns.items
          .sort(
            (a, b) =>
              new Date(
                b.status.conditions[
                  b.status.conditions.length - 1
                ].lastTransitionTime
              ) -
              new Date(
                a.status.conditions[
                  a.status.conditions.length - 1
                ].lastTransitionTime
              )
          )
          .reduce((result, pipelineRun) => {
            if (
              branches.indexOf(pipelineRun.metadata.labels.gitBranch) === -1
            ) {
              branches.push(pipelineRun.metadata.labels.gitBranch);
              const time = new Date(
                pipelineRun.status.conditions[
                  pipelineRun.status.conditions.length - 1
                ].lastTransitionTime
              );
              result.push({
                id: `${pipelineRun.metadata.labels.gitBranch}-branch`,
                branch: pipelineRun.metadata.labels.gitBranch,
                time: `${time.toLocaleDateString()} - ${time.toLocaleTimeString()}`,
                status:
                  pipelineRun.status.conditions[
                    pipelineRun.status.conditions.length - 1
                  ].reason
              });
            }
            return result;
          }, []);

        this.setState({
          rows,
          loading: false
        });
      })
      .catch(error => {
        error.response.text().then(text => {
          this.setState({
            error: text,
            rows: [],
            loading: false
          });
        });
      });
  }

  render() {
    const { close } = this.props;
    const { rows, loading, error } = this.state;

    const headers = [
      {
        key: 'branch',
        header: 'Branch'
      },
      {
        key: 'time',
        header: 'Last Build Time'
      },
      {
        key: 'status',
        header: 'Status'
      }
    ];

    return (
      <Modal
        open
        id="webhook-branches-modal"
        modalHeading="Latest PipelineRuns By Branch:"
        passiveModal
        onRequestClose={close}
      >
        {error && (
          <InlineNotification
            kind="error"
            subtitle={error}
            title="Error:"
            lowContrast
          />
        )}
        <div>
          <p>
            <strong>Webhook Name:</strong> {this.props.webhook.name}
          </p>
          <p>
            <strong>Repository:</strong> <a target="_blank" rel="noopener noreferrer" href={this.props.webhook.url}>{this.props.webhook.url}</a>
          </p>
          <p>
            <strong>Pipeline:</strong> {this.props.webhook.pipeline}
          </p>
          <p>
            <strong>Namespace:</strong> {this.props.webhook.namespace}
          </p>
        </div>
        <DataTable
          rows={rows}
          headers={headers}
          render={({ rows, headers, getHeaderProps, getRowProps }) => (
            <TableContainer>
              {loading ? (
                <DataTableSkeleton
                  rowCount={1}
                  columnCount={headers.length}
                  data-testid="loading-table"
                />
              ) : (
                <Table>
                  <TableHead>
                    <TableRow>
                      {headers.map(header => (
                        <TableHeader
                          key={header.id}
                          {...getHeaderProps({ header })}
                          isSortable
                          isSortHeader
                        >
                          {header.header}
                        </TableHeader>
                      ))}
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {rows.map(row => (
                      <TableRow {...getRowProps({ row })} key={row.id}>
                        {row.cells.map((cell, index) => (
                          <TableCell
                            className="cellText"
                            key={cell.id}
                            data-status={
                              index === row.cells.length - 1 ? cell.value : null
                            }
                          >
                            {cell.value}
                          </TableCell>
                        ))}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </TableContainer>
          )}
        />
        {rows.length === 0 && !loading && (
          <div className="noBranches">
            <p>
              Unable to identify any PipelineRuns initiated by this webhook.
            </p>
          </div>
        )}
      </Modal>
    );
  }
}
