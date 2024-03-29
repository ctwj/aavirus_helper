import React, { useEffect, useState } from 'react';
import { Layout, Nav, Button, Toast, Skeleton, Avatar } from '@douyinfe/semi-ui';
import { IconBell, IconHelpCircle, IconBytedanceLogo, IconHome, IconHistogram, IconLive, IconSetting, IconSemiLogo } from '@douyinfe/semi-icons';
import { TerminalOutput, TerminalInput } from 'react-terminal-ui'
import { useStore, observer } from './hooks/storeHook';

// 特征码定位
import SignatureCode from './page/signature_code/index'
// 帮助页
import Help from './page/help';
// 命令日志
import Log from './page/log';

import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'

import { CloseApp } from "../wailsjs/go/project/Project"


const App = () => {
    const { Header, Footer, Sider, Content } = Layout;
    const { appStore } = useStore()
    const [closing, setClosing] = useState(false)

    const [page, setPage] = useState('Home')

    // Nav 选中项
    const defaultSelectedKeys = [page]

    // Nav 选中事件
    const onSelect = (item) => {
        const { itemKey } = item
        setPage(itemKey)
    }

    // 处理go传递过来的消息
    useEffect(() => {
        EventsOn('message', (optionalData) => {
            appStore.addTerminalLineData(<TerminalOutput>{optionalData}</TerminalOutput>)
        })
        EventsOn('command', (optionalData) => {
            appStore.addTerminalLineData(<TerminalInput>{optionalData}</TerminalInput>)
        })
        EventsOn('progress', (optionalData) => {
            appStore.setProgress(optionalData)
        })
        return () => {
            EventsOff('message')
            EventsOff('command')
        }
    }, [])

    // 关闭已经打开的app
    const closeApp = () => {
        if (appStore.packing) {
            Toast.error({
                content: '打包中！请在空闲状态再关闭',
                duration: 3,
            })
            return;
        }
        setClosing(true)

        // 未反编译， 直接关闭
        if (!appStore.disassembled) {
            appStore.closeApp()
            setClosing(false)
            return
        }

        // 反编译了， 需要移除代码
        CloseApp(appStore.disassembleDir).then(result => { // 关闭app
            setClosing(false)
            appStore.closeApp()
        })
    }

    return (
        <Layout style={{ border: '1px solid var(--semi-color-border)', height: '100%' }}>
            <Sider style={{ backgroundColor: 'var(--semi-color-bg-1)' }}>
                <Nav
                    defaultSelectedKeys={defaultSelectedKeys}
                    onSelect={onSelect}
                    style={{ maxWidth: 220, height: '100%' }}
                    items={[
                        { itemKey: 'Home', text: '特征码定位', icon: <IconHome size="large" /> },
                        { itemKey: 'Log', text: '运行日志', icon: <IconHistogram size="large" /> },
                        { itemKey: 'Help', text: '使用说明', icon: <IconLive size="large" /> },
                        // { itemKey: 'Setting', text: '设置', icon: <IconSetting size="large" /> },
                    ]}
                    header={{
                        logo: <IconSemiLogo style={{ fontSize: 36 }} />,
                        text: 'Android 辅助去毒',
                    }}
                    footer={{
                        collapseButton: true,
                    }}
                    defaultIsCollapsed
                />
            </Sider>
            <Layout>
                <Header style={{ backgroundColor: 'var(--semi-color-bg-1)' }}>
                    <Nav
                        mode="horizontal"
                        footer={
                            <>
                                <Button
                                    theme="borderless"
                                    icon={<IconBell size="large" />}
                                    style={{
                                        color: 'var(--semi-color-text-2)',
                                        marginRight: '12px',
                                    }}
                                />
                                <Button
                                    theme="borderless"
                                    icon={<IconHelpCircle size="large" />}
                                    style={{
                                        color: 'var(--semi-color-text-2)',
                                        marginRight: '12px',
                                    }}
                                />
                                <Avatar color="orange" size="small">
                                    YJ
                                </Avatar>
                            </>
                        }
                    >
                    <span>{appStore?.path?.split('/')?. pop()}</span>
                    {appStore.path && <Button theme='borderless' type='primary' size='small' style={{ marginLeft: '12px' }}
                        onClick={closeApp}
                        disabled={closing}>关闭Apk</Button>}
                    </Nav>
                </Header>
                <Content
                    style={{
                        padding: '24px',
                        backgroundColor: 'var(--semi-color-bg-0)',
                        height: 'calc(100vh - 120px)',
                    }}
                >
                    {/* <Breadcrumb
                        style={{
                            marginBottom: '24px',
                        }}
                        routes={['首页', '当这个页面标题很长时需要省略', '上一页', '详情页']}
                    /> */}
                    <div
                        style={{
                            borderRadius: '10px',
                            border: '1px solid var(--semi-color-border)',
                            height: '100%',
                            boxSizing: 'border-box',
                        }}
                    >
                        { page === 'Home' && <SignatureCode /> /** 首页显示 特征码定位 */}
                        { page === 'Help' && <Help /> /** 首页显示 特征码定位 */}
                        { page === 'Log' && <Log /> /** 首页显示 特征码定位 */}
                        
                    </div>
                </Content>
                <Footer
                    style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        padding: '20px',
                        color: 'var(--semi-color-text-2)',
                        backgroundColor: 'rgba(var(--semi-grey-0), 1)',
                    }}
                >
                    <span
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                        }}
                    >
                        <IconBytedanceLogo size="large" style={{ marginRight: '8px' }} />
                        <span>Copyright © Ctwj. All Rights Reserved. </span>
                    </span>
                    <span>
                        <span>反馈建议</span>
                    </span>
                </Footer>
            </Layout>
        </Layout>
    );
}

export default observer(App)
