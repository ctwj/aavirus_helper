import { observer, useStore } from "../../hooks/storeHook";

import NoApk from "../../component/NoApk";
import Controller from "../../component/Controller";


// 处理 apk 文件打包页面
const File = () => {

    const { appStore } = useStore()
    return appStore.path ? <Controller /> : <NoApk />
}

export default observer(File)