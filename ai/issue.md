 Feature Clarification: AdminGroupBinding

That's what my understanding is:




Not sure what EntraID / OIDC security groups are, but the idea is that we want to be able to grant admin privileges to a whole group of users, rather than having to manually edit each individual user. This way, when a user is added or removed from the group in the identity provider (IdP), their admin status will automatically update in our system.
 

is this correct?



Curren How it is done ATM:


- k get users -A
-  k get user 3113fe4417b51774263f113407c327d3135edee3879fc5c1fad5b71a6877f45e -oyaml | grep -A20 admin




---
Investigate Issue: https://github.com/kubermatic/kubermatic/issues/14761


We need to implement this but before we do that, we need to investigate the issue and understand and how current implementation works. Investigate properly and help to undertsand what changes required to implement the new feature.