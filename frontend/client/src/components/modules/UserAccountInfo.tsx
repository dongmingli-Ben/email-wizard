import React from "react";
import { Box, Grid, IconButton, Typography } from "@mui/material";
import AddCircleOutlineIcon from "@mui/icons-material/AddCircleOutline";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";

type UserAccountInfoProps = {
  userName: string;
  userAccounts: { address: string; protocol: string }[];
  setAddAccount: (status: boolean) => void;
  setDeleteAccount: (address: string) => void;
};

const UserMailboxRow = (props: {
  address: string;
  protocol: string;
  setDeleteAccount: (address: string) => void;
}) => {
  return (
    <Grid
      item
      sx={{
        width: "100%",
      }}
    >
      <Box
        sx={{
          "&:hover": {
            backgroundColor: "primary.dark",
            opacity: [0.9, 0.8, 0.7],
            cursor: "default",
          },
          width: "100%",
          pl: "10%",
          pr: "5%",
          pt: 0.2,
          pb: 0.2,
          boxSizing: "border-box",
          display: "flex",
          alignItems: "center",
        }}
      >
        <Typography
          color="inherit"
          sx={{
            width: "100%",
          }}
          noWrap
        >
          {props.address}
        </Typography>
        <IconButton
          onClick={() => {
            props.setDeleteAccount(props.address);
          }}
          color="inherit"
        >
          <DeleteOutlineIcon></DeleteOutlineIcon>
        </IconButton>
      </Box>
    </Grid>
  );
};

const UserAccountInfo = (props: UserAccountInfoProps) => {
  let rows = props.userAccounts.map(
    (account: { address: string; protocol: string }, index: number) => {
      return (
        <UserMailboxRow
          address={account.address}
          protocol={account.protocol}
          setDeleteAccount={props.setDeleteAccount}
          key={index}
        />
      );
    }
  );
  return (
    <Grid
      container
      direction="column"
      justifyContent="flex-start"
      sx={{
        color: "common.white",
        mt: 3,
        width: "100%",
      }}
    >
      <Grid
        item
        sx={{
          width: "100%",
        }}
      >
        <Box
          sx={{
            "&:hover": {
              backgroundColor: "primary.dark",
              opacity: [0.9, 0.8, 0.7],
              cursor: "default",
            },
            width: "100%",
            pl: "10%",
            pr: "5%",
            pt: 1,
            pb: 1,
            boxSizing: "border-box",
          }}
        >
          <Typography
            color="inherit"
            variant="h4"
            noWrap
            sx={{
              width: "100%",
            }}
          >
            {props.userName}
          </Typography>
        </Box>
      </Grid>
      {props.userAccounts.length > 0 ? rows : <></>}
      <Grid item>
        <Box
          sx={{
            "&:hover": {
              backgroundColor: "primary.dark",
              opacity: [0.9, 0.8, 0.7],
              cursor: "default",
            },
            width: "100%",
            pl: "10%",
            pr: "5%",
            boxSizing: "border-box",
          }}
        >
          <IconButton
            onClick={() => {
              props.setAddAccount(true);
            }}
            sx={{
              color: "inherit",
              pl: 0,
              pb: 1,
            }}
          >
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
              }}
            >
              <AddCircleOutlineIcon color="inherit"></AddCircleOutlineIcon>
              <Typography
                color="inherit"
                sx={{
                  width: "100%",
                  fontWeight: "bold",
                  ml: 0.5,
                }}
                noWrap
                // gutterBottom
              >
                Add Mailbox
              </Typography>
            </Box>
          </IconButton>
        </Box>
      </Grid>
    </Grid>
  );
};

export default UserAccountInfo;
