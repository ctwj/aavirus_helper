import { observer, useStore } from "../../hooks/storeHook";

import NoApk from "../../component/NoApk";
import SmailComp from "../../component/Smail";


// 处理 apk 文件打包页面
const Smail = () => {

    const { appStore } = useStore()
    return appStore.path ? <SmailComp /> : <NoApk />
}

export default observer(Smail)