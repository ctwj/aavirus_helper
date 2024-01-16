import { observer } from "../hooks/storeHook";
import React, { useState } from "react";
import { Spin, Divider, Button, Typography } from '@douyinfe/semi-ui';
import TerminalUI from './Terminal';

import { Disassemble } from '../../wailsjs/go/project/Project'

import { useStore } from "../hooks/storeHook";

const ToolBar = () => {

    const { appStore } = useStore()
    const { disassembled } = appStore

    // 设置正在汇编的状态
    const [disassembling, setDisassembling] = useState()
    

    const  handleDisassemble = () => {
        setDisassembling(true)
        Disassemble(appStore.path).then(() => {
            setDisassembling(false)
        })
    }

    return (
        <React.Fragment>
            {!disassembled &&  
            <Button style={{ padding: '6px 24px',alignSelf: 'flex-start'  }} theme="solid" type="primary"
                onClick={handleDisassemble}
                disabled={disassembling}>
                {!disassembling ? '反编译' : '反编译中'}
            </Button>
            }
            {disassembling && <Spin />}
            {disassembled && <div></div>}
        </React.Fragment>
    )
}

const Controller = () => {

    const { appStore } = useStore()
    const { disassembled, apkInfo } = appStore
    const { Text } = Typography;
    

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <ToolBar style={{alignSelf: 'flex-start'}} />
            <Divider margin='12px'/>
            <div style={{ display: 'flex', flexDirection: 'column', flexShrink: '1', flexGrow: '1', overflow: 'auto' }}>
            <TerminalUI />
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