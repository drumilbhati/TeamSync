import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldSet,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { Card } from "@/components/ui/card";

const Login = () => {
  const [formData, setFormData] = useState({ email: "", password: "" });
  {
    /*const [error, setError] = useState(null);*/
  }
  const [token, setToken] = useState("");

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleLogin = async () => {
    try {
      const response = await fetch("/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });
      if (!response.ok) {
        throw new Error("Network response was not OK");
      }

      const data = await response.json();
      setToken(data.token || data);
      {
        /*setError(null); */
      }
      console.log({ token });
    } catch (error) {
      console.log(error);
      {
        /* setError(error.message); */
      }
    }
  };

  return (
    <div className="flex w-full min-h-screen justify-center items-center">
      <div className="w-full max-w-md">
        <Card className="bg-[#121212]">
          <FieldSet className="items-center justify-center">
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="username">Email</FieldLabel>
                <Input
                  id="email"
                  type="text"
                  placeholder="example@teamsync.com"
                  value={formData.email}
                  onChange={handleChange}
                />
                <FieldDescription>
                  Choose a unique email for your account.
                </FieldDescription>
              </Field>
              <Field>
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <FieldDescription>
                  Must be at least 8 characters long.
                </FieldDescription>
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  value={formData.password}
                  onChange={handleChange}
                />
              </Field>
            </FieldGroup>
          </FieldSet>
          <div className="flex justify-center">
            <Button onClick={handleLogin}>Login</Button>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Login;
