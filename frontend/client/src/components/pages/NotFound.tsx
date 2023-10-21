import React from "react";
import { Box, Button, Typography } from "@mui/material";
import Container from "@mui/material/Container";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";

// TODO remove, this demo shouldn't need to reset the theme.
const defaultTheme = createTheme();

export default function NotFound() {
  return (
    <ThemeProvider theme={defaultTheme}>
      <Container component="main" maxWidth="sm">
        <CssBaseline />
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flexDirection: "column",
            minHeight: "100vh",
          }}
        >
          <Typography variant="h1">404</Typography>
          <Typography variant="h6" sx={{ m: 2 }}>
            The page you’re looking for doesn’t exist.
          </Typography>
          <Button variant="contained" href="/">
            Back Home
          </Button>
        </Box>
      </Container>
    </ThemeProvider>
  );
}
