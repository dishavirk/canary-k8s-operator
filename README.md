# Canary K8s Operator
The Canary K8s Operator manages the deployment of Canary releases in our K8s cluster by leveraging Custom Resource Definitions (CRDs) and a Custom Controllers.
</br>
Recognizing a glaring gap between the abundance of theoretical articles on K8s Deployment Strategies 📚 and their practical application, I was inspired to make my own small contribution toward bridging this divide. This led me to spend a somewhat challenging weekend immersed in new concepts and topics, facing unexpected errors, and engaging in some endless debugging sessions.
I'm sharing my insights and lessons in the first part of my series, focusing on writing custom K8s operators using the Operator SDK.

## Description
Whether you're new to K8s or looking to expand your knowledge on Operator pattern and Canary deployments, please find the corresponding series on Medium for more context and deeper understanding that enhances the practical application of the concepts discussed.

[Part 1: Introduction to Operators and the Operator Pattern](https://medium.com/@disha.20.10/understanding-canary-deployments-in-kubernetes-part-1-introduction-to-operators-and-the-operator-0a0483d499a2)
Dive into the basics of K8s operators, understanding their purpose, and how they simplify complex operations tasks within a Kubernetes environment.

[Part 2: Implementing CRDs, Controllers, and Testing](https://medium.com/@disha.20.10/understanding-canary-deployments-in-kubernetes-part-2-implementing-crds-controllers-and-testing-3c3672edd99c)
Explore the practical aspects of developing a K8s operator, focusing on CRDs, writing controllers, and testing your operator.
