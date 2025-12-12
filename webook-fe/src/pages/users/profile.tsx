import { ProDescriptions } from '@ant-design/pro-components';
import React, { useState, useEffect } from 'react';
import { Button, message, Space } from 'antd';
import axios from "@/axios/axios";
import { useRouter } from 'next/router';

function Page() {
    let p: Profile = {Email: "", Phone: "", Nickname: "", Birthday:"", AboutMe: ""}
    const [data, setData] = useState<Profile>(p)
    const [isLoading, setLoading] = useState(false)
    const router = useRouter();

    useEffect(() => {
        setLoading(true)
        axios.get('/users/profile')
            .then((res) => res.data)
            .then((resp) => {
                setData(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)
                // 处理未登录
                if (err.response && err.response.status === 401) {
                    message.error('请先登录');
                    router.push('/users/login');
                } else {
                    message.error('获取用户信息失败');
                }
            })
    }, [])

    if (isLoading) return <p>Loading...</p>
    if (!data) return <p>No profile data</p>

    return (
        <>
            <Space style={{ marginBottom: 16 }}>
                <Button href={"/articles/list"}>返回主页</Button>
                <Button href={"/users/edit"} type={"primary"}>修改</Button>
            </Space>
            <ProDescriptions
                column={1}
                title="个人信息"
            >
                <ProDescriptions.Item label="昵称" valueType="text">
                    {data.Nickname}
                </ProDescriptions.Item>
                <ProDescriptions.Item valueType="text" label="邮箱">
                    {data.Email}
                </ProDescriptions.Item>
                <ProDescriptions.Item valueType="text" label="手机">
                    {data.Phone}
                </ProDescriptions.Item>
                <ProDescriptions.Item label="生日" valueType="date">
                    {data.Birthday}
                </ProDescriptions.Item>
                <ProDescriptions.Item valueType="text" label="关于我">
                    {data.AboutMe}
                </ProDescriptions.Item>
            </ProDescriptions>
        </>
    )
}

export default Page
