import React from "react";
import { observer } from "../hooks/storeHook";
import { useStore } from "../hooks/storeHook";

import { Tree, Toast } from '@douyinfe/semi-ui'


// type FileInfo struct {
// 	Label        string     `json:"label"`
// 	Value        string     `json:"value"`
// 	Key          string     `json:"key"`
// 	Name         string     `json:"name"`
// 	Size         int64      `json:"size"`
// 	IsDir        bool       `json:"isDir"`
// 	Path         string     `json:"path"`
// 	TotalSize    int64      `json:"totalSize,omitempty"`
// 	TotalFileNum int        `json:"totalFileNum,omitempty"`
// 	HumanSize    string     `json:"humanSize,omitempty"`
// 	Children     []FileInfo `json:"children,omitempty"`
// }

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


    const nodeStyle= {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
    };
    const handleRenderData = (list) => {
        list.array.forEach(element => {
            if (element.children) {
                handleRenderData(element)
            }

            if (element.isDir) { // 如果是 dir
                const label = element.label
                element.label = <ButtonGroup
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
                </ButtonGroup>
            } else { // 如果是文件

            }
        });
    }

    // 自定义渲染Tree节点
    const renderData = handleRenderData(fileListTreeData)

    const style = {
        width: '100%',
        height: '100%',
        border: '1px solid var(--semi-color-border)'
    };
    return (
        <div style={{ height: '100%' }}>
        <Tree
            treeData={renderData}
            directory
            showLine
            style={style}
        />
        </div>
    )
}

export default observer(FileList)