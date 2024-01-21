import { observer } from "../hooks/storeHook";
import React, { useEffect, useState } from "react";
import { SplitButtonGroup, Dropdown, Spin, RadioGroup, Radio, Button, Typography, Space } from '@douyinfe/semi-ui';
import { IconTreeTriangleDown } from '@douyinfe/semi-icons';
import XMLViewer from 'react-xml-viewer'

import { Disassemble, AndroidManifestInfo, BatchPack, OpenOutput } from '../../wailsjs/go/project/Project'

import Permission from "./Permission";

import { useStore } from "../hooks/storeHook";
import { useDisassemble } from "../hooks/disassemble";

// 打包按钮
const PackBtn = observer(() => {
    const { appStore } = useStore()
    const { disassembleApk, packApkWithDeleteDir } = useDisassemble()
    const { Text } = Typography

    // 打包
    const handleDir = (mode) => {
        packApkWithDeleteDir(mode)
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

    return <Space>
        <SplitButtonGroup>
            <Button theme="solid" type="primary">打包</Button>
            <Dropdown onVisibleChange={(v)=>setBtnVisible(v)} menu={menu} trigger="click" position="bottomRight">
                <Button style={btnVisible ? { background: 'var(--semi-color-primary-hover)', padding: '8px 4px' } : { padding: '8px 4px' }} theme="solid" type="primary" icon={<IconTreeTriangleDown />}></Button>
            </Dropdown>
        </SplitButtonGroup>
        <Button theme="solid" type="tertiary">全部打包</Button>
        <Button theme="solid" type="tertiary">全部交叉打包</Button>
    </Space>

})

// 工具栏
const ToolBar = observer(({setContent, setParseData, tab, setTab}) => {

    const { appStore } = useStore()
    const { disassembled, disassembleDir } = appStore
    const { disassembleApk, packApkWithDeleteDir } = useDisassemble()
    const { Text } = Typography;

    // 设置正在汇编的状态
    const [disassembling, setDisassembling] = useState()
    
    // 反编译 apk, 加载反编译结果的文件列表
    const  handleDisassemble = () => {
        setDisassembling(true)  // 开启编译中状态
        disassembleApk(appStore.path).then
        Disassemble(appStore.path).then(AndroidManifestInfo).then(result => {// 获取反编译后的文件列表
            const { content, parse } = result;
            setContent(content)
            setParseData(parse)
        }) 
    }

    // 打开output目录
    const openOutput = () => {
        OpenOutput()
    }

    useEffect(() => {
        // 加载完成， 如果 已经反汇编过， 在尝试获取一次文件列表
        if (appStore.disassembled) {
            AndroidManifestInfo(appStore.disassembleDir).then(result => {
                const { content, parse } = result;
                setContent(content)
                setParseData(parse)
            })
        }
        
    }, [])

    // 修改 Tabs
    const changeTabs = (e) => {
        console.log(e)
        setTab(e.target.value)
    }

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
                    <RadioGroup
                        onChange={changeTabs}
                        value={tab}
                        type='button'
                        style={{
                            display: 'flex',
                            width: 200,
                            justifyContent: 'center',
                            marginRight: '8px'
                        }}
                    >
                        <Radio value={'Permission'}>权限</Radio>
                        <Radio value={'Activity'}>活动</Radio>
                        <Radio value={'Permissions & Activities'}>综合</Radio>
                    </RadioGroup>
                    <PackBtn />
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
const Manifest = () => {

    const { appStore } = useStore()

    // AndroidManifest 的文本内容
    const [ content, setContent ] = useState("")
    const [ parseData, setParseData ] = useState({}) 

    const { UsesPermissions = [], Application = {} } = parseData ?? {};
    const { Activities = [], Services = [] } = Application ?? {};

    // console.log(UsesPermissions, Activities, Services)

    const [ tab, setTab ] = useState('Permission')

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <ToolBar setContent={setContent} setParseData={setParseData} tab={tab} setTab={setTab} />
            <div style={{ display: 'flex', flexDirection: 'column', flexShrink: '1', flexGrow: '1', overflow: 'auto' }}>
                { tab === 'Permission' && <Permission permissions={UsesPermissions} /> }
                { tab === 'Activity' && <Permission data={UsesPermissions} /> }
                { tab === 'Permissions & Activities' && <Permission data={UsesPermissions} /> }
                {/* <XMLViewer  xml={content} /> */}
            </div>
        </div>
    )
}

export default observer(Manifest)