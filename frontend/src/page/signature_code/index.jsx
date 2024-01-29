import { memo } from 'react';
import { useStore, observer } from "../../hooks/storeHook";

import File from './file'
import Manifest from './manifest'
import Smail from "./smail";
import Upload from "./upload";

const FileMemo = memo(() => {
    // File 组件的渲染逻辑
    return <File />;
  });
  
  const ManifestMemo = memo(() => {
    // Manifest 组件的渲染逻辑
    return <Manifest />;
  });
  


// 特征码定位
const SignatureCode = () => {

    const { appStore } = useStore()

    if (appStore.func === 'file') { // 文件处理
        return <FileMemo />
    }

    if (appStore.func === 'manifest') { // 配置处理
        return <ManifestMemo />
    }

    if (appStore.func === 'upload') {
        return <Upload />
    }

    return <Smail /> // 代码处理

}

export default observer(SignatureCode)