import React, { useState, useEffect } from 'react';
import axios from "@/axios/axios";
import {useSearchParams} from "next/navigation";
import {Content} from "antd/es/layout/layout";
import {Typography} from "antd";
import Title from "antd/es/typography/Title";
import {ProLayout} from "@ant-design/pro-components";

function Page() {
    const [data, setData] = useState<Article>()
    const [isLoading, setLoading] = useState(false)
    const params = useSearchParams()
    const artID = params?.get("id")
    useEffect(() => {
        setLoading(true)
        axios.get('/articles/pub/'+artID)
            .then((res) => res.data)
            .then((data) => {
                setData(data.data)
                setLoading(false)
            })
    }, [artID])

    if (isLoading) return <p>Loading...</p>
    if (!data) return <p>No data</p>

    return (
        <ProLayout pure={true}>
            <Typography>
                <Title>
                    {data.title}
                </Title>
                <Content dangerouslySetInnerHTML={{__html: data.content}}>
                </Content>
            </Typography>
        </ProLayout>
    )
}

export default Page