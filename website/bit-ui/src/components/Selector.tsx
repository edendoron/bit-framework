import React, {FC} from 'react';
import {FormControl, FormHelperText, InputLabel, MenuItem, Select} from "@material-ui/core";

interface SelectorProps {
    menuItems: string[],
    queryType: string,
    onChange: (event: React.ChangeEvent<{ value: unknown }>) => void,
}

export const Selector: FC<SelectorProps> = ({menuItems, queryType, onChange}) => {

    const renderMenuItems = () => {
        return menuItems.map(item => <MenuItem value={item}>{item}</MenuItem>)
    }
    return (
        <FormControl>
            <Select
                value={queryType}
                onChange={onChange}
            >
                {renderMenuItems()}
            </Select>
        </FormControl>
        )
}