import React, {useEffect, useState} from 'react';
import Axios from 'axios';
import './App.css';
import {Box, Card, CardContent, CardHeader} from '@material-ui/core'
import {Selector} from "./components/Selector";

const STORAGE_DATA_READ_URL = 'http://localhost:8082/data/read'
const queryTypes = ['Reports', 'BIT Status', 'Config Failures']
const userGroups = ['group1', 'group2', 'group3', 'group4', 'groupRafael', 'TemperatureCelsius group', 'group general', 'groupField']

export const App = () => {
    const [queryType, setQueryType] = useState('');
    const [userGroup, setUserGroup] = useState('');
    const [data, setData] = useState();

    const changeQueryType = (event: React.ChangeEvent<{ value: unknown }>) => {
        setQueryType(event.target.value as string)
    }

    const changeUserGroup = (event: React.ChangeEvent<{ value: unknown }>) => {
        setUserGroup(event.target.value as string)
    }

    useEffect(() => {
        Axios.get(STORAGE_DATA_READ_URL + '?config_failures').then((res) => {
            setData(res.data)
        })
    }, [queryType]);

    const renderData = () => {
        switch (queryType){
            case 'Reports':
                break;
            case 'BIT Status':
                break;
            case 'Config Files':
                // fetchFilteringRules()
        }
        return <div/>;
    }

    // const renderSearch = () => {
    //
    // }

    return (
    <Box bgcolor="#373A36" minHeight="100vh" textAlign="center">
        <Box bgcolor="#D48166" minHeight="100vh" marginRight="150px" marginLeft="150px">
            <Card>
                <CardHeader title="BIT Framework Query System"/>
                <CardContent>
                    <Selector menuItems={userGroups} queryType={userGroup} onChange={changeUserGroup}/>
                    {userGroup !== '' && <Selector menuItems={queryTypes} queryType={queryType} onChange={changeQueryType}/>}

                    {renderData()}
                </CardContent>
            </Card>
        </Box>
    </Box>
    );
}

export default App;
