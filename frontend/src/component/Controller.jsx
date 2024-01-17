import { observer } from "../hooks/storeHook";
import React, { useState } from "react";
import { Spin, Divider, Button, Typography } from '@douyinfe/semi-ui';
import TerminalUI from './Terminal';
import FileList from "./FileList";

import { Disassemble, FileList as GetFileList } from '../../wailsjs/go/project/Project'

import { useStore } from "../hooks/storeHook";

const ToolBar = () => {

    const { appStore } = useStore()
    const { disassembled } = appStore

    // 设置正在汇编的状态
    const [disassembling, setDisassembling] = useState()
    
    // 反编译 apk, 加载反编译结果的文件列表
    const  handleDisassemble = () => {
        setDisassembling(true)  // 开启编译中状态
        Disassemble(appStore.path).then((result) => { // 反编译完成后， 设置状态
            setDisassembling(false)
            appStore.setDisassembled(true)
            appStore.setDisassembleDir(result.outdir)
            return result.outdir
        }).then(GetFileList).then(result => {// 获取反编译后的文件列表
            appStore.setDisassembleFileList(result)
        }) 
    }

    return (
        <div style={{height: '36px', padding: "8px 16px"}}>
            {!disassembled &&  
            <Button style={{ padding: '6px 24px',alignSelf: 'flex-start'  }} theme="solid" type="primary"
                onClick={handleDisassemble}
                disabled={disassembling}>
                {!disassembling ? '反编译' : '反编译中'}
            </Button>
            }
            {disassembling && <Spin />}
            {disassembled && <div></div>}
        </div>
    )
}

const Controller = () => {

    const { appStore } = useStore()
    const { disassembled, apkInfo } = appStore
    const { Text } = Typography;
    

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <ToolBar />
            {/* <Divider margin='12px'/> */}
            <div style={{ display: 'flex', flexDirection: 'column', flexShrink: '1', flexGrow: '1', overflow: 'auto' }}>
            { !disassembled && <TerminalUI />}
            { disassembled && <FileList />}
                {/* {
                    apkInfo.map((item, index) => {
                        return <Text key={index} style={{ whiteSpace: 'pre-wrap'}}>{item}</Text>
                    })
                } */}
            </div>
        </div>
    )
}

export default observer(Controller)