'use client';
import {EditOutlined, EyeOutlined} from '@ant-design/icons';
import {ProLayout, ProList} from '@ant-design/pro-components';
import {Button, Tabs, Tag} from 'antd';
import React, {useEffect, useMemo, useState} from 'react';
import axios from "@/axios/axios";
import router from "next/router";
import dynamic from "next/dynamic";

enum ArticleStatus {
    Unknown = 0,
    Draft = 1,
    Withdraw = 2,
    Published = 3,
    Deleted = 4,
}

interface ArticleItem {
    id: number
    title: string
    status: ArticleStatus
    abstract: string
    author?: string
}

const IconButton = ({icon, text, onClick}: { icon: any, text: string, onClick: any }) => (
    <Button onClick={onClick} type={"default"}>
        {React.createElement(icon, {style: {marginInlineEnd: 8}})}
        {text}
    </Button>
);

const statusTag = (status: ArticleStatus) => {
    switch (status) {
        case ArticleStatus.Draft:
            return <Tag color="processing">草稿</Tag>
        case ArticleStatus.Withdraw:
            return <Tag color="warning">已撤回</Tag>
        case ArticleStatus.Published:
            return <Tag color="success">已发布</Tag>
        default:
            return <Tag>未知</Tag>
    }
}

const ArticleList = () => {
    const [mine, setMine] = useState<Array<ArticleItem>>([])
    const [publicArticles, setPublicArticles] = useState<Array<ArticleItem>>([])
    const [loadingMine, setLoadingMine] = useState<boolean>(false)
    const [loadingPub, setLoadingPub] = useState<boolean>(false)
    const [activeTab, setActiveTab] = useState<string>('published')

    useEffect(() => {
        setLoadingMine(true)
        axios.post('/articles/list', {
            offset: 0,
            limit: 100,
        }).then((res) => res.data)
            .then((data) => {
                setMine(data.data || [])
            })
            .finally(() => setLoadingMine(false))
    }, [])

    useEffect(() => {
        setLoadingPub(true)
        axios.post('/articles/pub/list', {
            offset: 0,
            limit: 100,
        }).then((res) => res.data)
            .then((data) => {
                setPublicArticles(data.data || [])
            })
            .finally(() => setLoadingPub(false))
    }, [])

    const myDrafts = useMemo(
        () => mine.filter(a => a.status === ArticleStatus.Draft || a.status === ArticleStatus.Withdraw),
        [mine],
    )
    const myPublished = useMemo(
        () => mine.filter(a => a.status === ArticleStatus.Published),
        [mine],
    )

    const renderList = (data: ArticleItem[], loading: boolean, isMine: boolean) => (
        <ProList<ArticleItem>
            toolBarRender={() => {
                return isMine ? [
                    <Button key="create" type="primary" href={"/articles/edit"}>
                        写作
                    </Button>,
                ] : [];
            }}
            itemLayout="vertical"
            rowKey="id"
            headerTitle={isMine ? "我的文章" : "公开文章"}
            loading={loading}
            dataSource={data}
            metas={{
                title: {
                    dataIndex: "title"
                },
                description: {
                    render: (_, record) => (
                        <>
                            {statusTag(record.status)}
                            {record.author ? <Tag style={{marginLeft: 8}}>{record.author}</Tag> : null}
                        </>
                    )
                },
                actions: {
                    render: (_, row) => {
                        if (isMine) {
                            const actions: React.ReactNode[] = [
                                <IconButton
                                    icon={EditOutlined}
                                    text="编辑"
                                    onClick={() => router.push("/articles/edit?id=" + row.id.toString())}
                                    key="edit"
                                />,
                            ]
                            if (row.status === ArticleStatus.Published) {
                                actions.push(
                                    <IconButton
                                        icon={EyeOutlined}
                                        text="查看"
                                        onClick={() => router.push("/articles/view?id=" + row.id.toString())}
                                        key="view"
                                    />
                                )
                                actions.push(
                                    <Button
                                        danger
                                        type="default"
                                        key="withdraw"
                                        onClick={() => {
                                            axios.post('/articles/withdraw', {id: row.id})
                                                .then(res => res.data)
                                                .then(res => {
                                                    if (res.code === 0) {
                                                        // move to drafts
                                                        setMine(prev => prev.map(it => it.id === row.id ? {...it, status: ArticleStatus.Withdraw} : it))
                                                    }
                                                })
                                        }}
                                    >
                                        撤回
                                    </Button>
                                )
                            }
                            return actions
                        }
                        return [
                            <IconButton
                                icon={EyeOutlined}
                                text="查看"
                                onClick={() => router.push("/articles/view?id=" + row.id.toString())}
                                key="view"
                            />,
                        ]
                    }
                },
                extra: {
                    render: () => (
                        <img
                            width={240}
                            alt="cover"
                            src="https://gw.alipayobjects.com/zos/rmsportal/mqaQswcyDLcXyDKnZfES.png"
                        />
                    ),
                },
                content: {
                    render: (_, record) => (
                        <div dangerouslySetInnerHTML={{__html: record.abstract}}></div>
                    )
                },
            }}
        />
    )

    return (
        <ProLayout title={"创作中心"}>
            <div style={{display: "flex", justifyContent: "flex-end", marginBottom: 12}}>
                <Button key="profile" href={"/users/profile"} type="default">
                    个人信息
                </Button>
            </div>
            <Tabs
                activeKey={activeTab}
                onChange={setActiveTab}
                items={[
                    {
                        key: 'published',
                        label: '我的已发布',
                        children: renderList(myPublished, loadingMine, true),
                    },
                    {
                        key: 'drafts',
                        label: '草稿 / 撤回',
                        children: renderList(myDrafts, loadingMine, true),
                    },
                    {
                        key: 'public',
                        label: '公开文章',
                        children: renderList(publicArticles, loadingPub, false),
                    },
                ]}
            />
        </ProLayout>
    );
};

// 关闭 SSR 避免 ProLayout 在首屏渲染时的水合差异
export default dynamic(() => Promise.resolve(ArticleList), { ssr: false });
