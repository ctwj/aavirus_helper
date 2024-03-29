import React from "react";
import { observer } from "../hooks/storeHook";
import { useStore } from "../hooks/storeHook";

import { Typography, ButtonGroup, Button, Tree, Toast } from '@douyinfe/semi-ui'

const nodeStyle= {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center'
};
const TreeItemNode = (props) => {
    const node = props.node
    const { name, humanSize } = node;

    const { appStore } = useStore()
    const showManiFest = () => {
        appStore.setFunc('manifest')
    }

    const { Text } = Typography

    return <div style={nodeStyle}>
        <span>{name}</span>
        <div style={{ display: 'flex', alignItems: 'center'}}>
            <Text style={{marginRight: '8px'}}>{humanSize}</Text>
            {name === 'AndroidManifest.xml' && <Button
                    theme="borderless"
                    type='secondary'
                    onClick={showManiFest}
                >分析</Button>}
            {/* <ButtonGroup
                size="small"
                theme="borderless"
                >
                <Button
                    onClick={e => {
                        Toast.info(opts);
                        e.stopPropagation();
                    }}
                >提示</Button>
                <Button>点击</Button>
            </ButtonGroup> */}
        </div>
    </div>
}

/**
 * FileList
 * @param {*} props 
 * @returns 
 */
const FileList = (props) => {

    const { appStore } = useStore()
    const { fileListTreeData } = appStore

    const handleTerminalInput = (terminalInput) => {
    }


    const treeSelectChange = e => {
        appStore.setSelFiles(e)
    } 
    
    const handleRenderData = (list) => {
        if (!list.length) {
            return list;
        }
        list.map(element => {
            if (element.children) {
                element.children = handleRenderData(element.children)
            }

            const label = element.label
            if (element.isDir) { // 如果是 dir
                element.label = <TreeItemNode node={element} />
            } else { // 如果是文件
                element.label = <TreeItemNode node={element} />
            }
            return element;
        });
        return list;
    }

    // 自定义渲染Tree节点
    const renderData = handleRenderData(JSON.parse(JSON.stringify(fileListTreeData)))

    const style = {
        width: '100%',
        height: '100%',
        border: '1px solid var(--semi-color-border)'
    };
    return (
        <div style={{ height: '100%' }}>
        <Tree
            treeData={renderData}
            defaultValue={appStore.selFiles}
            directory
            showLine
            multiple
            checkRelation='unRelated'
            style={style}
            onChange={treeSelectChange}
        />
        </div>
    )
}

export default observer(FileList)