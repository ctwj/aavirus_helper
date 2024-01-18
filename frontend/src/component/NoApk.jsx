import React from 'react';
import { IconBell, IconHelpCircle, IconBytedanceLogo, IconHome, IconHistogram, IconLive, IconSetting, IconSemiLogo } from '@douyinfe/semi-icons';
import { Empty, Button } from '@douyinfe/semi-ui';
import { IllustrationNoContent, IllustrationNoContentDark } from '@douyinfe/semi-illustrations';

import { useStore } from '../hooks/storeHook'

import { OpenFile, GetApkInfo } from '../../wailsjs/go/project/Project'
import { LogPrint } from '../../wailsjs/runtime/runtime'

function NoApk() {

    const { appStore } =  useStore()

    const openFile = () => {
        OpenFile().then((result) => {
            const apkPath = result.file
            appStore.setApkPath(apkPath)
            GetApkInfo(apkPath).then(result => {
                LogPrint(JSON.stringify(result))
                appStore.setApkInfo(result.info)
            })
        })
    }

    return (
        <div style={{ height: '80%', display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
            <Empty
                image={<IllustrationNoContent style={{ width: 150, height: 150 }} />}
                darkModeImage={<IllustrationNoContentDark style={{ width: 150, height: 150 }} />}
                title="第一步从选择一个apk开始"
                description="请选择需要去毒的apk文件"
            >
                <div style={{ display: 'flex', justifyContent: 'center' }}>
                    <Button style={{ padding: '6px 24px' }} theme="solid" type="primary"
                        onClick={openFile}>
                        选择
                    </Button>
                </div>
            </Empty>
        </div>
    );
};

export default NoApk
