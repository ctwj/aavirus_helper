import { makeObservable, observable, action, computed } from "mobx"

export class App {
    id = Math.random()
    path = ""
    apkInfo = []

    // 命令行输出结果
    terminalLineData = []

    disassembled = false
    disassembleDir = ""
    disassembleFileList = {}

    finished = false

    constructor() {
        makeObservable(this, {
            path: observable,
            apkInfo: observable,
            finished: observable,
            disassembled: observable,        // 是否已经反汇编
            disassembleDir: observable,      // 反汇编文件目录
            disassembleFileList: observable, // 反汇编文件夹的文件列表
            terminalLineData: observable,
            unfinishedTodoCount: computed,
            fileListTreeData: computed,
            toggle: action,
            setDisassembled: action,        // 设置 disassembled 状态
            setDisassembleDir: action,      // 设置反编译后的文件夹
            setDisassembleFileList: action, // 设置反编译后的文件列表
            setApkPath: action,             // 设置应用的路径
            addTerminalLineData: action,    // 向命令行输出数据中，添加一行
            setApkInfo: action,             // 设置 apk 的基本信息
        })
    }

    get unfinishedTodoCount() {
        return this.path.length;
    }

    get fileListTreeData() {
        return [this.disassembleFileList];
    }

    setApkPath (path) {
        this.path = path;
    }

    setApkInfo (info) {
        this.apkInfo = info;
    }

    setDisassembled (disassembled) {
        this.disassembled = disassembled
    }

    setDisassembleDir (disassembleDir) {
        this.disassembleDir = disassembleDir
    }

    setDisassembleFileList (fileList) {
        this.disassembleFileList = fileList
    }

    toggle() {
        this.finished = !this.finished
    }

    addTerminalLineData (line) {
        this.terminalLineData = [...this.terminalLineData, line]
    }
}