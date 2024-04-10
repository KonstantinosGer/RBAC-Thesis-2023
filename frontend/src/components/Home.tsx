import {Card, Col, Row} from 'antd';
import React from 'react';
import {Link} from "react-router-dom";

const InfoCard: React.FC<{
    title: string;
    index: number;
    desc: string;
    href: string;
}> = ({title, href, index, desc}) => {
    return (
        <div
            style={{
                backgroundColor: '#FFFFFF',
                boxShadow: '0 2px 4px 0 rgba(35,49,128,0.02), 0 4px 8px 0 rgba(49,69,179,0.02)',
                borderRadius: '8px',
                fontSize: '14px',
                color: 'rgba(0,0,0,0.65)',
                textAlign: 'justify',
                lineHeight: ' 22px',
                padding: '16px 19px',
                // flex: 1,
            }}
        >
            <div
                style={{
                    display: 'flex',
                    gap: '4px',
                    alignItems: 'center',
                }}
            >
                <div
                    style={{
                        width: 48,
                        height: 48,
                        lineHeight: '22px',
                        backgroundSize: '100%',
                        textAlign: 'center',
                        padding: '8px 16px 16px 12px',
                        color: '#FFF',
                        fontWeight: 'bold',
                        backgroundImage:
                            "url('https://gw.alipayobjects.com/zos/bmw-prod/daaf8d50-8e6d-4251-905d-676a24ddfa12.svg')",
                    }}
                >
                    {index}
                </div>
                <div
                    style={{
                        fontSize: '16px',
                        color: 'rgba(0, 0, 0, 0.85)',
                        paddingBottom: 8,
                    }}
                >
                    {title}
                </div>
            </div>
            <div
                style={{
                    fontSize: '14px',
                    color: 'rgba(0,0,0,0.65)',
                    textAlign: 'justify',
                    lineHeight: '22px',
                    marginBottom: 8,
                }}
            >
                {desc}
            </div>
            <Link to={href}>
                Go {'>'}
            </Link>
            {/*<a href={href} target="_blank" rel="noreferrer">*/}
            {/*    Link {'>'}*/}
            {/*</a>*/}
        </div>
    );
};

const Home: React.FC = () => {
    return (
        <Card
            style={{
                borderRadius: 8,
            }}
            bodyStyle={{
                backgroundImage:
                    'radial-gradient(circle at 97% 10%, #EBF2FF 0%, #F5F8FF 28%, #EBF1FF 124%)',
            }}
        >
            <div
                style={{
                    backgroundPosition: '100% -8%',
                    backgroundRepeat: 'no-repeat',
                    backgroundSize: '274px auto',
                    backgroundImage:
                        "url('https://gw.alipayobjects.com/mdn/rms_a9745b/afts/img/A*BuFmQqsB2iAAAAAAAAAAAAAAARQnAQ')",
                }}
            >
                <div
                    style={{
                        fontSize: '20px',
                        color: '#1A1A1A',
                    }}
                >
                    RBAC
                </div>
                <p
                    style={{
                        fontSize: '14px',
                        color: 'rgba(0,0,0,0.65)',
                        lineHeight: '22px',
                        marginTop: 16,
                        marginBottom: 32,
                        width: '65%',
                    }}
                >
                    Role Based Access Control System
                </p>
                <div
                    style={{
                        gap: 16,
                    }}
                >
                    <Row justify={"space-around"} gutter={80} style={{width: '94%', marginLeft: '3%'}}>
                        <Col span={8}>
                        {/*<Col>*/}
                            <InfoCard
                                index={1}
                                title="Users"
                                href="/users"
                                desc="All firebase users"
                            />
                        </Col>
                        <Col span={8}>
                        {/*<Col>*/}
                            <InfoCard
                                index={2}
                                title="Employees"
                                href="/employees"
                                desc="All DM members"
                            />
                        </Col>
                        <Col span={8}>
                        {/*<Col>*/}
                            <InfoCard
                                index={3}
                                title="Customers"
                                href="/customers"
                                desc="All DM clients"
                            />
                        </Col>
                    </Row>
                    <br/><br/>
                    <Row justify={"space-around"} style={{width: '90%', verticalAlign: "middle", marginLeft: '5%'}}>
                        <Col span={8}>
                        {/*<Col>*/}
                            <InfoCard
                                index={4}
                                title="Roles"
                                href="/roles"
                                desc="All members with their roles"
                            />
                        </Col>
                        <Col span={8}>
                        {/*<Col>*/}
                            <InfoCard
                                index={5}
                                title="Permissions"
                                href="/permissions"
                                desc="All roles with their system permissions"
                            />
                        </Col>
                    </Row>
                </div>
            </div>
        </Card>
    );
};

export default Home;
