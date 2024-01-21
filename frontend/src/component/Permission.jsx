import { useStore, observer } from "../hooks/storeHook"
import { Table } from '@douyinfe/semi-ui';

const Permission = ({permissions}) => {

    let i = 0
    const data = permissions?.map(item => {
        item.key = i++;
        item.name = item.Name;
        item.description = item.Name;
        return item;
    })

    const columns = [
        {
            title: '权限',
            dataIndex: 'name',
        },
        {
            title: '说明',
            dataIndex: 'description',
        },
    ];

    const rowSelection = {
        getCheckboxProps: record => ({
            disabled: record.name === '设计文档', // Column configuration not to be checked
            name: record.name,
        }),
        onSelect: (record, selected) => {
            console.log(`select row: ${selected}`, record);
        },
        onSelectAll: (selected, selectedRows) => {
            console.log(`select all rows: ${selected}`, selectedRows);
        },
        onChange: (selectedRowKeys, selectedRows) => {
            console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
        },
    };

    return <Table style={{width: '100%'}}
                sticky={{ top: 0 }} 
                columns={columns} 
                rowSelection={rowSelection} 
                dataSource={data} 
                pagination={false} />;
}

export default observer(Permission)