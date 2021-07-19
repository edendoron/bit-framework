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
    TableRow, TableSortLabel,
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

const compareDate = (failure1: failureObject, failure2: failureObject) => {
    const Date1 = new Date(failure1.timestamp.seconds * 1000).getTime();
    const Date2 = new Date(failure2.timestamp.seconds * 1000).getTime();
    return Date1 > Date2 ? 1 : -1;
}


const FailureRow = (failure: failureObject) => {
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

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
                <TableCell align="center">{failure.failure_data.bit_type.map((type) => type + ', ')}</TableCell>
                <TableCell align="center">{failure.count}</TableCell>
                <TableCell align="center">{failure.failure_data.severity}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{paddingBottom: 0, paddingTop: 0}} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box margin={1}>
                            <div><span
                                className={classes.tableHeader}>Description:</span> "{failure.failure_data.description}"
                            </div>

                            <div><span
                                className={classes.tableHeader}>Additional Info:</span> "{failure.failure_data.additional_info}"
                            </div>

                            <div><span className={classes.tableHeader}>Purpose:</span> "{failure.failure_data.purpose}"
                            </div>

                            <div><span
                                className={classes.tableHeader}>Operator Failure:</span> "{failure.failure_data.operator_failure}"
                            </div>

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

const StatusRow = (props: { status: statusObject, index: number, time: number }) => {
    const {status, index} = props;
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    const timestamp = new Date(props.time).toLocaleString()

    return (
        <React.Fragment>
            <TableRow className={classes.root}>
                <TableCell>
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/>}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">{index}</TableCell>
                <TableCell align="center">{timestamp}</TableCell>
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

const compareStatusDate = (status1: statusObject, status2: statusObject) => {
    const Date1 = new Date(status1.failures[0].timestamp.seconds * 1000).getTime();
    const Date2 = new Date(status2.failures[0].timestamp.seconds * 1000).getTime();
    return Date1 > Date2 ? 1 : -1;
}

const compareNumber = (status1: statusObject, status2: statusObject) => {
    return status1.failures.length > status2.failures.length ? 1 : -1;
}

function stableSort<T>(array: T[], comparator: (a: T, b: T) => number, direction: string) {
    const stabilizedThis = array.map((el, index) => [el, index] as [T, number]);
    stabilizedThis.sort((a, b) => {
        const order = direction === "desc" ? comparator(a[0], b[0]) : -comparator(a[0], b[0]);
        if (order !== 0) return order;
        return a[1] - b[1];
    });
    return stabilizedThis.map((el) => el[0]);
}


const getComparator = (key: string): ((a: statusObject, b: statusObject) => number) => {
    switch (key) {
        case "Date":
            return compareStatusDate
        case "# of Failures":
            return compareNumber
        default:
            return compareStatusDate
    }
}

const useHeaderRowStyles = makeStyles({
    root: {
        '& > *': {
            borderBottom: 'set',
            backgroundColor: '#b3a48e',
        },
    },
    text: {
        fontWeight: 'bold',
    }
});
export const StatusTable: FC<StatusTableProps> = ({data}) => {
    const [orderByField, setOrderField] = React.useState("No.")
    const [orderByDirection, setOrderDirection] = React.useState<"asc" | "desc">("asc")

    const classes = useHeaderRowStyles();

    if (!data || data.length === 0) return <div>No Statuses Found.</div>

    return (
        <TableContainer component={Paper}>
            <Table aria-label="collapsible table">
                <TableHead className={classes.root}>
                    <TableRow>
                        <TableCell/>
                        <TableCell className={classes.text}>No.</TableCell>
                        <TableCell
                            className={classes.text}
                            align="center"
                            sortDirection={false}>
                            <TableSortLabel
                                active={orderByField === "Date"}
                                direction={orderByDirection}
                                onClick={() => {
                                    if (orderByDirection === "asc")
                                        setOrderDirection("desc")
                                    else
                                        setOrderDirection("asc")
                                    setOrderField("Date")
                                }}>
                                Timestamp
                            </TableSortLabel>
                        </TableCell>
                        <TableCell
                            className={classes.text}
                            align="center"
                            sortDirection={false}>
                            <TableSortLabel
                                active={orderByField === "# of failures"}
                                direction={orderByDirection}
                                onClick={() => {
                                    if (orderByDirection === "asc")
                                        setOrderDirection("desc")
                                    else
                                        setOrderDirection("asc")
                                    setOrderField("# of failures")
                                }}>
                                # of failures
                            </TableSortLabel>
                        </TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {stableSort(data, getComparator(orderByField), orderByDirection).map((status, index) =>
                    <StatusRow status={status} index={index} time={status.failures[0].timestamp.seconds * 1000}/>
                    )}
                </TableBody>
            </Table>
        </TableContainer>
    )
}