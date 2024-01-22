import { useStore, observer } from "../hooks/storeHook"
import { Table } from '@douyinfe/semi-ui';


// props: {permissions: '权限列表', setSelectedItem: '设置选中项'}
const Permission = ({permissions, setSelectedItem}) => {

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
        // getCheckboxProps: record => ({
        //     // disabled: record.name === '设计文档', // Column configuration not to be checked
        //     // name: record.name,
        // }),
        // onSelect: (record, selected) => {
        //     console.log(`select row: ${selected}`, record);
        // },
        onSelectAll: (selected, selectedRows) => {
            setSelectedItem(selectedRows);
        },
        onChange: (selectedRowKeys, selectedRows) => {
            setSelectedItem(selectedRows);
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