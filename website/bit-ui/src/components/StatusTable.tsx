import React, {FC} from 'react';
import {
    Box,
    Collapse,
    IconButton,
    makeStyles,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography
} from "@material-ui/core";
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

interface failureObject {
    failure_data: {
        unit_name: string,
        test_name: string,
        test_id: number,
        bit_type: Array<string>,
        description: string,
        additional_info: string,
        purpose: string,
        severity: number,
        operator_failure: Array<string>,
        line_replacent_units: Array<string>,
        field_replacemnt_units: Array<string>,
    }
    timestamp: { nanos: number, seconds: number },
    count: number,
}

interface statusObject {
    failures: Array<failureObject>,
}

interface StatusTableProps {
    data: Array<statusObject>,

}

const useRowStyles = makeStyles({
    root: {
        '& > *': {
            borderBottom: 'set',
            backgroundColor: '#D48166',
        },
    },
    tableHeader: {
        textDecorationLine: 'underline',
        fontWeight: 'bold',
    }
});

const compareDate = (report1: failureObject, report2: failureObject) => {
    const Date1 = new Date(report1.timestamp.seconds * 1000).getTime();
    const Date2 = new Date(report2.timestamp.seconds * 1000).getTime();
    return Date1 > Date2 ? 1 : -1;
}


const FailureRow = (failure: failureObject) => {
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    const timestamp = new Date(failure.timestamp.seconds * 1000).toLocaleString()

    return (
        <React.Fragment>
            <TableRow style={{backgroundColor: 'beige'}}>
                <TableCell>
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/>}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">{failure.failure_data.test_id}</TableCell>
                <TableCell align="center">{failure.failure_data?.test_name}</TableCell>
                <TableCell align="center">{timestamp}</TableCell>
                <TableCell align="center">{failure.failure_data.bit_type.map((type) => type + ', ')}</TableCell>
                <TableCell align="center">{failure.count}</TableCell>
                <TableCell align="center">{failure.failure_data.severity}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{paddingBottom: 0, paddingTop: 0}} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box margin={1}>
                            <div><span className={classes.tableHeader}>Description:</span> "{failure.failure_data.description}"</div>

                            <div><span className={classes.tableHeader}>Additional Info:</span> "{failure.failure_data.additional_info}"</div>

                            <div><span className={classes.tableHeader}>Purpose:</span> "{failure.failure_data.purpose}"</div>

                            <div><span className={classes.tableHeader}>Operator Failure:</span> "{failure.failure_data.operator_failure}"</div>

                            <div><span className={classes.tableHeader}>Line Replacement Units:</span>
                                "{failure.failure_data.line_replacent_units.map((unit) => unit + ', ')}"
                            </div>

                            <div><span className={classes.tableHeader}>Field Replacement Units:</span>
                                "{failure.failure_data.field_replacemnt_units.map((unit) => unit + ', ')}"
                            </div>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    )
}

const StatusRow = (props: { status: statusObject, index: number }) => {
    const {status, index} = props;
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    return (
        <React.Fragment>
            <TableRow className={classes.root}>
                <TableCell>
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/>}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">{index}</TableCell>
                <TableCell align="center">{status.failures.length}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{paddingBottom: 0, paddingTop: 0}} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box margin={1}>
                            <Typography variant="h6" gutterBottom component="div">
                                Failures
                            </Typography>
                            <TableContainer component={Paper}>
                                <Table aria-label="collapsible table">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell/>
                                            <TableCell>ID</TableCell>
                                            <TableCell align="center">Name</TableCell>
                                            <TableCell align="center">Timestamp</TableCell>
                                            <TableCell align="center">BIT Type</TableCell>
                                            <TableCell align="center">Failure Count</TableCell>
                                            <TableCell align="center">Severity</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {status.failures.sort(compareDate).map((failure) => (
                                            FailureRow(failure)
                                        ))}

                                    </TableBody>
                                </Table>
                            </TableContainer>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    )
}

export const StatusTable: FC<StatusTableProps> = ({data}) => {
    if(data.length === 0) return <div>No Statuses Found.</div>

    return (
        <TableContainer component={Paper}>
            <Table aria-label="collapsible table">
                <TableHead>
                    <TableRow>
                        <TableCell/>
                        <TableCell>No.</TableCell>
                        <TableCell align="center"># of Failures</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {
                        data.map((status, index) => (
                            <StatusRow status={status} index={index}/>
                        ))}
                </TableBody>
            </Table>
        </TableContainer>
    )
}