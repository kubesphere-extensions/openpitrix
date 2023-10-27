import React, { useRef, useState } from 'react';
import styled from 'styled-components';
import { Link, useParams } from 'react-router-dom';
import { Banner, BannerTip, Button, Field, notify } from '@kubed/components';

import {
  Icon,
  Column,
  TableRef,
  DataTable,
  StatusIndicator,
  useItemActions,
  useTableActions,
  useBatchActions,
  DeleteConfirmModal,
  useListQueryParams,
} from '@ks-console/shared';

import RepoManagementModal from './RepoManagementModal';
import { getRepoUrl, useReposDeleteMutation } from '../../stores';
import type { RepoData } from '../../types/repo';

const AddButton = styled(Button)`
  min-width: 96px;
`;

function Repos(): JSX.Element {
  const { workspace = '' } = useParams();
  const repoListUrl = getRepoUrl({ workspace });
  const tableRef = useRef<TableRef>();
  const [modalType, setModalType] = useState<string>('');
  const [selectedRows, setSelectedRows] = useState<RepoData[]>();
  const { mutateAsync, isLoading } = useReposDeleteMutation(workspace);
  const tableParameters = {
    order: 'create_time',
    status: 'active',
  };
  const renderItemActions = useItemActions<RepoData>({
    authKey: 'app-repos',
    params: { workspace },
    actions: [
      {
        key: 'edit',
        icon: <Icon name="pen" />,
        text: t('EDIT_INFORMATION'),
        action: 'edit',
        onClick: (_, record) => {
          setSelectedRows([record]);
          setModalType('edit');
        },
      },
      {
        key: 'delete',
        icon: <Icon name="trash" />,
        text: t('DELETE'),
        action: 'delete',
        onClick: (_, record) => {
          setSelectedRows([record]);
          setModalType('delete');
        },
      },
    ],
  });
  const renderBatchActions = useBatchActions({
    authKey: 'app-repos',
    params: { workspace },
    actions: [
      {
        key: 'delete',
        text: t('DELETE'),
        action: 'delete',
        onClick: () => {
          const selectedFlatRows = tableRef?.current?.getSelectedFlatRows() || [];
          setSelectedRows(selectedFlatRows as unknown as RepoData);
          setModalType('delete');
        },
        props: { color: 'error' },
      },
    ],
  });
  const renderTableActions = useTableActions({
    authKey: 'app-repos',
    params: { workspace },
    actions: [
      {
        key: 'create',
        text: t('ADD'),
        action: 'manage',
        onClick: () => setModalType('create'),
        props: {
          color: 'secondary',
          shadow: true,
        },
      },
    ],
  });
  const columns: Column<RepoData>[] = [
    {
      title: t('NAME'),
      field: 'spec.name',
      width: '25%',
      searchable: true,
      render: (name, record) => (
        <Field
          value={<Link to={record.metadata.name}>{name}</Link>}
          label={name || '-'}
          avatar={<Icon name="catalog" size={40} />}
        />
      ),
    },
    {
      title: t('STATUS'),
      field: 'status.state',
      canHide: true,
      width: '15%',
      render: (status = 'syncing') => (
        <StatusIndicator type={status}>
          {t(`APP_REPO_STATUS_${status.toUpperCase()}`)}
        </StatusIndicator>
      ),
    },
    {
      title: t('URL'),
      field: 'spec.url',
    },
    {
      id: 'more',
      title: '',
      width: 20,
      render: renderItemActions,
    },
  ];

  function serverDataFormatter(serverData: any) {
    return {
      ...serverData,
      totalItems: serverData.total_count,
    };
  }

  function transformRequestParams(params: Record<string, any>): Record<string, any> {
    const { parameters, pageIndex, filters } = params;
    const keyword = filters?.[0]?.value;
    const formattedParams = useListQueryParams({
      ...parameters,
      page: pageIndex + 1,
    });

    if (!keyword) {
      return formattedParams;
    }

    return {
      ...formattedParams,
      conditions: formattedParams.conditions + `,keyword=${keyword}`,
    };
  }

  function closeModal(): void {
    setModalType('');
    setSelectedRows(undefined);
  }

  function handleManageOk(): void {
    tableRef.current?.refetch();
    closeModal();
  }

  async function handleRepoDelete(): Promise<void> {
    const reposId: string[] = selectedRows?.map(item => item.metadata.name) || [];
    await mutateAsync(reposId);
    notify.success(t('DELETED_SUCCESSFUL'));
    closeModal();
    tableRef.current?.refetch();
  }

  return (
    <>
      <Banner
        className="mb12"
        icon={<Icon name="catalog" />}
        title={t('APP_REPO')}
        description={t('APP_REPO_DESC')}
      >
        <BannerTip title={t('HOW_TO_USE_APP_REPO_Q')} key="develop">
          {t('HOW_TO_USE_APP_REPO_A')}
        </BannerTip>
      </Banner>
      <DataTable
        ref={tableRef}
        rowKey="metadata.uid"
        tableName="APP_REPOSITORY"
        url={repoListUrl}
        simpleSearch
        parameters={tableParameters}
        transformRequestParams={transformRequestParams}
        columns={columns}
        useStorageState={false}
        placeholder={t('SEARCH_BY_NAME')}
        toolbarRight={renderTableActions()}
        batchActions={renderBatchActions()}
        format={item => ({ ...item, workspace })}
        serverDataFormat={serverDataFormatter}
        emptyOptions={{
          withoutTable: true,
          createButton: (
            <AddButton color="secondary" onClick={() => setModalType('create')}>
              {t('ADD')}
            </AddButton>
          ),
          title: t('NO_APP_REPO_FOUND'),
          image: <Icon name="catalog" size={48} />,
          description: t('APP_REPOSITORY_EMPTY_DESC'),
        }}
      />
      {['create', 'edit'].includes(modalType) && (
        <RepoManagementModal
          visible={true}
          onCancel={closeModal}
          onOk={handleManageOk}
          detail={selectedRows?.[0]}
        />
      )}
      {modalType === 'delete' && (
        <DeleteConfirmModal
          visible={true}
          type="APP_REPOSITORY"
          resource={selectedRows?.map((item: any) => item.name)}
          onOk={handleRepoDelete}
          onCancel={closeModal}
          confirmLoading={isLoading}
        />
      )}
    </>
  );
}

export default Repos;