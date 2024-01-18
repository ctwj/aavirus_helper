import { observer } from "../hooks/storeHook";
import React, { useEffect, useState } from "react";
import { Spin, Toast, Divider, Button, Typography } from '@douyinfe/semi-ui';
import TerminalUI from './Terminal';
import FileList from "./FileList";

import { Disassemble, FileList as GetFileList, BatchPack, OpenOutput } from '../../wailsjs/go/project/Project'

import { useStore } from "../hooks/storeHook";

const ToolBar = observer(() => {

    const { appStore } = useStore()
    const { disassembled, disassembleDir } = appStore
    const { Text } = Typography;

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

    // 打包
    const handleDir = () => {
        appStore.setPacking(true)
        BatchPack(disassembleDir, appStore.selFiles).then(result => {
            appStore.setPacking(false)
            Toast.info({
                content: '打包完成',
                duration: 3,
            })
            OpenOutput()
        })
    }

    // 打开output目录
    const openOutput = () => {
        OpenOutput()
    }

    useEffect(() => {
        // 加载完成， 如果 已经反汇编过， 在尝试获取一次文件列表
        if (appStore.disassembled) {
            GetFileList(appStore.disassembleDir).then(result => {
                appStore.setDisassembleFileList(result)
            })
        }
        
    }, [])

    return (
        <div style={{height: '36px', padding: "8px 16px", display: "flex", justifyContent: "space-between", alignItems: "center"}}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-start'}}>
                {!disassembled &&  
                <Button style={{ padding: '6px 24px',alignSelf: 'flex-start'  }} theme="solid" type="primary"
                    onClick={handleDisassemble}
                    disabled={disassembling}>
                    {!disassembling ? '反编译' : '反编译中'}
                </Button>
                }
                {disassembled && <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                    <Button style={{ padding: '6px 24px',alignSelf: 'flex-start'  }} theme="solid" type="primary"
                        onClick={handleDir}
                        disabled={appStore.packing}>
                        打包 {appStore.packing ? '中' : ''}
                    </Button>
                </div>}
                
                {(disassembling || appStore.packing) && <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                    <Spin style={{marginLeft: '8px'}} />
                    <Text type="tertiary" style={{marginLeft: '8px'}} >{ appStore.progress }</Text>
                </div>}
            </div>
            
            <Button style={{ padding: '6px 24px',alignSelf: 'flex-start', marginLeft: 'auto'  }} type="secondary"
                theme='borderless'
                onClick={openOutput}>
                打开文件夹
            </Button>
        </div>
    )
})

const Controller = () => {

    const { appStore } = useStore()
    const { disassembled, apkInfo } = appStore

    // 主面板显示的内容
    const MainPanel = observer(() => {
        if (!disassembled || appStore.packing) {
            return <TerminalUI />
        }
        return <FileList />
    })
    

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <ToolBar />
            <div style={{ display: 'flex', flexDirection: 'column', flexShrink: '1', flexGrow: '1', overflow: 'auto' }}>
             <MainPanel />
            </div>
        </div>
    )
}

export default observer(Controller)