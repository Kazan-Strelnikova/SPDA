import React, { useState, useEffect, useContext } from 'react';
import { Box, Button, TextField, Typography, Container, Paper, Link } from '@mui/material';
import variables from '../variables.module.scss'; 
import { UserContext } from '../contexts/UserContext';
import { Exception } from 'sass';
import axios from 'axios';
import { postLoginUser } from '../http/post-login-user';
import { useNavigate } from 'react-router-dom';

export const LogInPage: React.FC = () => {
    const navigate = useNavigate()
    const userContext = useContext(UserContext);
    if (!userContext) {
        throw new Error("Log must be used within a UserProvider");
    }
    const { user, setUser } = userContext;

    const [formData, setFormData] = useState({
        email: '',
        password: ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value
        });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        (() => {(async function (){
            try {
                const u = await postLoginUser(formData.email, formData.password)
                setUser(u)
                navigate("/")
            } catch (err: any) {
                console.log("login err", err)
            }
            
        }())}) ()


    };
    
    useEffect(() => {
        document.body.style.overflow = 'hidden';
        return () => {
            document.body.style.overflow = 'auto';
        };
    }, []);


    return (
        <Box sx={{ 
            minHeight: '70vh', 
            display: 'flex',
            justifyContent: 'center', 
            alignItems: 'center',
        }}>
            <Container maxWidth="xs" sx={{ my: 4 }}>
                <Box sx={{ 
                    display: 'flex', 
                    flexDirection: 'column',
                    alignItems: 'center',
                    mb: 2
                }}>
                </Box>
                
                <Paper 
                    elevation={3} 
                    sx={{ 
                        p: 4, 
                        borderRadius: 2,
                        width: '100%',
                        backgroundColor: 'white',
                        boxShadow: '0px 4px 20px rgba(0, 0, 0, 0.05)'
                    }}
                >
                    <Box component="form" onSubmit={handleSubmit} noValidate>
                        <TextField
                            margin="normal"
                            // required
                            fullWidth
                            id="email"
                            label="Email"
                            name="email"
                            autoComplete="email"
                            autoFocus
                            value={formData.email}
                            onChange={handleChange}
                            sx={{ mb: 2 }}
                            variant="outlined"
                        />

                        <TextField
                            margin="normal"
                            // required
                            fullWidth
                            name="password"
                            label="Пароль"
                            type="password"
                            id="password"
                            autoComplete="current-password"
                            value={formData.password}
                            onChange={handleChange}
                            sx={{ mb: 3 }}
                        />

                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{ 
                                py: 1.5, 
                                backgroundColor: variables.primary,
                                '&:hover': {
                                    backgroundColor:variables.dark
                                },
                                mb: 2
                            }}
                        >
                            Войти
                        </Button>

                        <Box sx={{ textAlign: 'center', mt: 1 }}>
                            <Typography variant="body2">
                                Нет аккаунта?{' '}
                                <Link 
                                    href="/signup" 
                                    variant="body2" 
                                    sx={{ 
                                        color: variables.dark,
                                        textDecoration: 'none',
                                        fontWeight: 500,
                                        '&:hover': {
                                            textDecoration: 'underline', 
                                        }
                                    }}
                                >
                                    Зарегистрироваться
                                </Link>
                            </Typography>
                        </Box>
                    </Box>
                </Paper>
            </Container>
        </Box>
    );
};