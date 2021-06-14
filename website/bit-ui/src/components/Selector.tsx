import React, {FC} from 'react';
import {FormControl, FormHelperText, InputLabel, MenuItem, Select} from "@material-ui/core";

interface SelectorProps {
    queryType: string,
    onChange: (event: React.ChangeEvent<{ value: unknown }>) => void,
}

export const Selector: FC<SelectorProps> = ({queryType, onChange}) => {

    return (
        <FormControl>
            <InputLabel>Type</InputLabel>
            <Select
                value={queryType}
                onChange={onChange}
            >
                <MenuItem value={"Reports"}>Reports</MenuItem>
                <MenuItem value={"BIT Status"}>Bit Status</MenuItem>
                <MenuItem value={"Config Files"}>Config Files</MenuItem>
            </Select>
            <FormHelperText>Select query type</FormHelperText>
        </FormControl>
        )
}