import { observer } from "../hooks/storeHook";
import React, { useEffect, useState } from "react";
import { SplitButtonGroup, Dropdown, Spin, Toast, Divider, Button, Typography } from '@douyinfe/semi-ui';
import { IconTreeTriangleDown } from '@douyinfe/semi-icons';

import TerminalUI from './Terminal';
import FileList from "./FileList";

import { Disassemble, FileList as GetFileList, BatchPack, OpenOutput } from '../../wailsjs/go/project/Project'

import { useStore } from "../hooks/storeHook";

// 打包按钮
const PackBtn = observer(() => {
    const { appStore } = useStore()
    const { Text } = Typography

    // 打包
    const handleDir = (mode) => {
        if (appStore.packing) {
            Toast.info({
                content: '请等待任务完成后再试',
                duration: 3,
            })
            return;
        }

        appStore.setPacking(true)
        BatchPack(appStore.disassembleDir, appStore.selFiles, mode).then(result => {
            appStore.setPacking(false)
            Toast.info({
                content: '打包完成',
                duration: 3,
            })
            OpenOutput()
        })
    }

    const menuGenerator = (selectItems) => {
        const menu = [
            { node: 'item', name: '单独模式', tag: 'single', disabled: true, onClick: () => handleDir('single') },
            { node: 'item', name: '组合模式', tag: 'group', disabled: true, onClick: () =>  handleDir('group') },
            { node: 'item', name: '交叉模式', tag: 'cross', disabled: true, onClick: () =>  handleDir('cross') },
        ]
        if (selectItems.length === 0) {
            return menu
        }

        // AndroidManifest.xml apktool.yml
        const items = selectItems.filter(item => !['AndroidManifest.xml', 'apktool.yml'].some(file => item.includes(file)))
        const length = items.length;

        return menu.map(item => {
            item.disabled = false
            if (item.tag === 'single') {
                item.name = <Text>单独模式 <Text type="tertiary">{length}</Text></Text>
            }
            if (item.tag === 'group') {
                if (length < 2) {
                    item.disabled = true
                }
                item.name = <Text>组合模式 <Text type="tertiary">1</Text></Text>
            }
            if (item.tag === 'cross') {
                if (length < 3) {
                    item.disabled = true
                    item.name = `交叉模式`
                } else {
                    item.name = <Text>交叉模式 <Text type="tertiary">{length * (length - 1) / 2}</Text></Text>
                }
            }
            return item
        })
    }
    const menu = menuGenerator(appStore.selFiles)

    const [btnVisible, setBtnVisible] = useState(false)

    return <SplitButtonGroup>
        <Button theme="solid" type="primary">打包</Button>
        <Dropdown onVisibleChange={(v)=>setBtnVisible(v)} menu={menu} trigger="click" position="bottomRight">
            <Button style={btnVisible ? { background: 'var(--semi-color-primary-hover)', padding: '8px 4px' } : { padding: '8px 4px' }} theme="solid" type="primary" icon={<IconTreeTriangleDown />}></Button>
        </Dropdown>
    </SplitButtonGroup>
})

// 工具栏
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
                    <PackBtn />
                    {/* <Button style={{ padding: '6px 24px',alignSelf: 'flex-start'  }} theme="solid" type="primary"
                        onClick={handleDir}
                        disabled={appStore.packing}>
                        打包 {appStore.packing ? '中' : ''}
                    </Button> */}
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

// 
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