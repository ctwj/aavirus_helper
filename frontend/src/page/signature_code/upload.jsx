import React, { useEffect, useState } from "react";
import QRCode from "qrcode.react"
import { 
    Card, ButtonGroup, InputNumber, Spin, 
    Switch, Popover, Button, Typography, Space, Toast
} from '@douyinfe/semi-ui';
import { IconLink } from '@douyinfe/semi-icons';
import XMLViewer from 'react-xml-viewer'

import { Disassemble, OpenOutput } from '../../../wailsjs/go/project/Project'
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime'


import { useStore, observer } from "../../hooks/storeHook";
import { StartServer, StopServer, GetStatus, GetIps } from "../../../wailsjs/go/upload/Web"

/**
 * 工具栏, 
 */
const ToolBar = () => {
    const { appStore } = useStore()

    // 打开output目录
    const openOutput = () => {
        OpenOutput()
    }

    // 显示文件列表页
    const backToController = () => {
        const old = appStore.oldFunc
        appStore.setFunc(old, old)
    }

    return (
        <div style={{height: '36px', padding: "8px 16px", display: "flex", justifyContent: "space-between", alignItems: "center"}}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-start', flexShrink: 1}}>
            </div>
            <ButtonGroup size="small" style={{ padding: '6px 24px',alignSelf: 'flex-start', marginLeft: 'auto', flexShrink: 0  }}>
                <Button type="secondary"
                    onClick={openOutput}>文件夹</Button>
                <Button type="secondary"
                    onClick={backToController}>返回</Button>
            </ButtonGroup>
        </div>
    )
}

// 
const Upload = () => {
    const { Text } = Typography;

    const [ips, setIps] = useState([])

    const [webStatus, setWebStatus] = useState(false)
    const [adbStatus, setAdbStatus] = useState(false)
    const webPort = localStorage.getItem('web_port');
    const [port, setPort] = useState({
        val: webPort || 1234,
        status: 'default',
    })
    
    // 修改端口
    const  changePort = (value, e) => {
        const valid = /\d+/.test(value) 
        if (!valid || value < 80 || value > 65535) {
            setPort({val: value, status: 'error'})
        } else {
            setPort({val: value, status: 'default'})
            localStorage.setItem('web_port', value)
        }
    }

    // 网站开关
    const webSwitch = (checked) => {
        serverCtl(checked)
    }

    const [serverChanging, setServerChanging] = useState(false)

    // 服务器控制
    const serverCtl = (status) => {
        setServerChanging(true)
        if (status) {
            StartServer(parseInt(port.val, 10)).then(result => {
                console.log(result)
                setWebStatus(result.status)
            }).finally (() => {
                console.log('serverChanging', serverChanging)
                setServerChanging(false)
            })
        } else {
            StopServer().then(result => {
                console.log(result)
            }).finally (() => {
                console.log('serverChanging', serverChanging)
                setServerChanging(false)
            })
        }
    }   

    // 初始化，读取 webserver 状态， 本机ip
    useEffect(() => {
        const statusCheck = () => {
            GetStatus(parseInt(port.val, 10)).then(result => {
                setWebStatus(result.status)
                console.log(result.status)
            })
        }
        const handler = setInterval(statusCheck, 2000)
        statusCheck()
        GetIps().then(setIps)
        return () => {
            clearInterval(handler)
        }
    }, [])

    const IpsRender = () => {
        return <Space>
            {ips.map((ip, index) => (
                <Popover
                    key={index}
                    content={
                        <QRCode value="http://{ip}:{port.val}" />
                    }
                >
                    <Text link icon={<IconLink />} underline
                        onClick={() => BrowserOpenURL(`http://${ip}:${port.val}`)}>http://{ip}:{port.val}</Text>
                </Popover>
            ))}
        </Space>
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <ToolBar  />
            <Space vertical style={{alignItems: 'stretch'}}>
                <Card 
                    title='adb 快捷上传' 
                    headerExtraContent={
                        <Text link>
                            更多
                        </Text>
                    }
                >
                    快捷上传，用于一键将打包apk上传至手机中 /sdcard/Download 目录中
                </Card>
                <Spin spinning={false} tip="loading...">
                <Card 
                    title='Output目录Web服务' 
                    style={{  }}
                    headerExtraContent={
                        <Switch checked={webStatus} onChange={webSwitch} aria-label="a switch for web server"></Switch>
                    }
                >
                    <Space vertical style={{alignItems: 'flex-start'}}>
                        <Text>启动服务后， 局域网内手机可以快捷访问本机ip地址，来下载apk，测试浏览器是否报毒</Text>
                        <div>
                            <InputNumber prefix="服务端口" 
                                value={port.val}
                                validateStatus={port.status}
                                disabled={webStatus}
                                max={65535}
                                min={80}
                                onChange={changePort}></InputNumber> <Text type="tertiary">端口范围在 80-65535 之间</Text>
                        </div>
                        <div>
                           { webStatus ? <IpsRender /> : <Text type="tertiary">服务未启动</Text>}
                        </div>
                    </Space>
                </Card>
                </Spin>
            </Space>
        </div>
    )
}

export default observer(Upload)