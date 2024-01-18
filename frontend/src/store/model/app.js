import { makeObservable, observable, action, computed } from "mobx"
import { makePersistable } from 'mobx-persist-store'

export class App {
    id = Math.random()
    
    // 是否反汇编了
    disassembled = false

    // 打包状态
    packing = false

    // 打包进度
    progress = ""

    // apk 路径
    path = ""

    // 反编译后代码路径
    disassembleDir = ""

    // 应用信息
    apkInfo = []

    // 命令行输出结果
    terminalLineData = []

    // 反编译目录文件列表
    disassembleFileList = {}

    // FileList 中选中的文件
    selFiles = []

    


    constructor() {
        makeObservable(this, {
            path: observable,
            apkInfo: observable,
            disassembled: observable,        // 是否已经反汇编
            disassembleDir: observable,      // 反汇编文件目录
            disassembleFileList: observable, // 反汇编文件夹的文件列表
            selFiles: observable,            // 选择的文件列表
            packing: observable,             // 是否正在打包
            progress: observable,            // 打包进度
            terminalLineData: observable,
            fileListTreeData: computed,
            setDisassembled: action,        // 设置 disassembled 状态
            setDisassembleDir: action,      // 设置反编译后的文件夹
            setDisassembleFileList: action, // 设置反编译后的文件列表
            setApkPath: action,             // 设置应用的路径
            addTerminalLineData: action,    // 向命令行输出数据中，添加一行
            setApkInfo: action,             // 设置 apk 的基本信息
            setSelFiles: action,            // 设置选择的数据
            setPacking: action,             // 设置是否正在打包
            setProgress: action,            // 设置打包进度
            closeApp: action,               // 关闭app
        })
        makePersistable(
            this, 
            { 
                name: 'SampleStore', 
                properties: [
                    'path',
                    'disassembled',
                    'disassembleDir',
                    {
                        key: 'disassembleFileList',
                        serialize: (value) => JSON.stringify(value),
                        deserialize: (value) => JSON.parse(value),
                    },
                ], 
                storage: window.localStorage
            }
        );

    }

    get isHydrated() {
        return isHydrated(this);
    }

    async getStoredData() {
        return getPersistedStore(this);
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

    setSelFiles (files) {
        this.selFiles = files
    }

    setPacking (val) {
        this.packing = val
    }

    setProgress (val) {
        this.progress = val
    }

    addTerminalLineData (line) {
        this.terminalLineData = [...this.terminalLineData, line]
    }

    closeApp () {
        this.path = ""
        this.packing = false
        this.disassembled = false
        this.apkInfo = []
        this.disassembleDir = ""
        this.terminalLineData = []
        this.disassembleFileList = {}
        this.selFiles = []
    }
}