import { makeObservable, observable, action, computed } from "mobx"

export class App {
    id = Math.random()
    path = ""
    apkInfo = []

    // 命令行输出结果
    terminalLineData = []

    disassembled = false

    finished = false

    constructor() {
        makeObservable(this, {
            path: observable,
            apkInfo: observable,
            finished: observable,
            disassembled: observable,
            terminalLineData: observable,
            unfinishedTodoCount: computed,
            toggle: action,
            setApkPath: action,             // 设置应用的路径
            addTerminalLineData: action,    // 向命令行输出数据中，添加一行
            setApkInfo: action,             // 设置 apk 的基本信息
        })
    }

    get unfinishedTodoCount() {
        return this.path.length;
    }

    setApkPath (path) {
        this.path = path;
    }

    setApkInfo (info) {
        this.apkInfo = info;
    }

    toggle() {
        this.finished = !this.finished
    }

    addTerminalLineData (line) {
        this.terminalLineData = [...this.terminalLineData, line]
    }
}